package jsonz

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Reconstruct is an example of how to reconstruct a JSON from Parse().
func Reconstruct(in []byte) ([]byte, error) {
	b := bytes.Buffer{}
	b.Grow(len(in))

	var lastTokenType TokenType
	for item, err := range Parse(in) {
		if err != nil {
			return nil, err
		}
		if ShouldAddComma(lastTokenType, item.Token.Type()) {
			b.WriteByte(',')
		}
		if item.Key.IsValue() {
			b.Write(item.Key.Raw())
			b.WriteByte(':')
		}
		b.Write(item.Token.Raw())
		lastTokenType = item.Token.Type()
	}
	return b.Bytes(), nil
}

// Reformat is an example of how to reconstruct a JSON from Parse(), and format with indentation.
func Reformat(in []byte, prefix, indent string) ([]byte, error) {
	b := bytes.Buffer{}
	b.Grow(len(in))

	var lastToken TokenType
	for item, err := range Parse(in) {
		if err != nil {
			return nil, err
		}
		if ShouldAddComma(lastToken, item.Token.Type()) {
			b.WriteByte(',')
		}
		if lastToken != 0 {
			b.WriteByte('\n')
		}
		b.WriteString(prefix)
		for range item.Level {
			b.WriteString(indent)
		}
		if item.Key.IsValue() {
			b.Write(item.Key.Raw())
			b.WriteString(": ")
		}
		b.Write(item.Token.Raw())
		lastToken = item.Token.Type()
	}
	return b.Bytes(), nil
}

type Builder struct {
	zBuf
	indent string
	prefix string

	lastTok TokenType
	level   int
	stack   []TokenType // array or object
	err     error

	lastIndent  bool
	lastNewline bool
}

// NewBuilder creates a new Builder. It's optional to set the prefix and indent. A zero Builder is valid.
func NewBuilder(prefix, indent string) *Builder {
	return &Builder{
		indent: indent,
		prefix: prefix,
	}
}

// Bytes returns the bytes of the builder with an error if any.
func (b *Builder) Bytes() ([]byte, error) {
	return b.zBuf, b.err
}

// AddRaw adds a key and token to the builder. It will add a comma if needed.
func (b *Builder) AddRaw(key, token RawToken) {
	switch {
	case token.IsOpen():
		b.WriteNewline(token.Type())
		b.WriteIndent()
		b.writeKey(key)
		b.writeByte(byte(token.Type()))
		b.setLastToken(token.Type())
		b.stack = append(b.stack, token.Type())
		b.level++

	case token.IsClose():
		if key.Type() != 0 {
			b.addErrorf("unexpected key(%s) before close token(%s)", key, token.Type())
			return
		}
		if b.level <= 0 {
			b.addErrorf("unexpected close token(%s)", token.Type())
			return
		}
		b.level--
		b.stack = b.stack[:len(b.stack)-1]
		b.WriteNewline(token.Type())
		b.WriteIndent()
		b.writeByte(byte(token.Type()))
		b.setLastToken(token.Type())

	case token.IsValue():
		b.WriteNewline(token.Type())
		b.WriteIndent()
		b.writeKey(key)
		b.write(token.Raw())
		b.setLastToken(token.Type())

	default:
		b.addErrorf("unexpected token(%s)", token)
	}
}

// AddRaw adds a key and token to the builder. It will add a comma if needed.
func (b *Builder) AddToken(key string, token RawToken) {
	switch {
	case token.IsOpen():
		b.WriteNewline(token.Type())
		b.WriteIndent()
		b.writeKeyString(key)
		b.writeByte(byte(token.Type()))
		b.setLastToken(token.Type())
		b.stack = append(b.stack, token.Type())
		b.level++

	case token.IsClose():
		if key != "" {
			b.addErrorf("unexpected key(%s) before close token(%s)", key, token.Type())
			return
		}
		if b.level <= 0 {
			b.addErrorf("unexpected close token(%s)", token.Type())
			return
		}
		b.level--
		b.stack = b.stack[:len(b.stack)-1]
		b.WriteNewline(token.Type())
		b.WriteIndent()
		b.writeByte(byte(token.Type()))
		b.setLastToken(token.Type())

	case token.IsValue():
		b.WriteNewline(token.Type())
		b.WriteIndent()
		b.writeKeyString(key)
		b.write(token.Raw())
		b.setLastToken(token.Type())

	default:
		b.addErrorf("unexpected token(%s)", token)
	}
}

// Add adds a key and value to the builder. It will add a comma if needed.
func (b *Builder) Add(key string, value any) {
	var tokType TokenType
	var rawTok RawToken

	switch v := value.(type) {
	case TokenType: // open or close token
		tokType = v

	case RawToken:
		rawTok = v
		tokType = rawTok.typ

	case []byte:
		var err error
		rawTok, err = NewRawToken(v)
		if err != nil {
			b.addErrorf(err.Error())
		}
		tokType = rawTok.typ

	default:
		// string | number | boolean | null
	}

	switch {
	case tokType.IsOpen():
		b.WriteNewline(tokType)
		b.WriteIndent()
		b.writeKeyString(key)
		b.writeByte(byte(tokType))
		b.setLastToken(tokType)
		b.stack = append(b.stack, tokType)
		b.level++

	case tokType.IsClose():
		if key != "" {
			b.addErrorf("unexpected key(%s) before close token(%s)", key, tokType)
			return
		}
		if b.level <= 0 {
			b.addErrorf("unexpected close token(%s)", tokType)
			return
		}
		b.level--
		b.stack = b.stack[:len(b.stack)-1]
		b.WriteNewline(tokType)
		b.WriteIndent()
		b.writeByte(byte(tokType))
		b.setLastToken(tokType)

	default:
		b.WriteNewline(TokenString) // assume value token
		b.WriteIndent()
		b.writeKeyString(key)

		if rawTok.typ != 0 {
			if !rawTok.IsValue() {
				b.addErrorf("unexpected token(%s) as value", rawTok)
				return
			}
			b.write(rawTok.raw)
			b.setLastToken(tokType)
		} else {
			b.writeValue(value)
		}
	}
}

// ShouldAddComma returns true if a comma should be added before the next token.
func (b *Builder) ShouldAddComma(next TokenType) bool {
	return ShouldAddComma(b.lastTok, next)
}

func (b *Builder) WriteComma(next TokenType) bool {
	ok := b.ShouldAddComma(next) || next == 0
	if ok {
		b.writeByte(',')
		b.setLastToken(TokenComma)
	}
	return ok
}

func (b *Builder) writeValue(value any) {
	b.setLastToken(TokenNumber) // default to number, can be overridden
	switch v := value.(type) {
	case nil:
		b.zBuf = append(b.zBuf, "null"...)
		b.setLastToken(TokenNull)
	case bool:
		if v {
			b.zBuf = append(b.zBuf, "true"...)
			b.setLastToken(TokenTrue)
		} else {
			b.zBuf = append(b.zBuf, "false"...)
			b.setLastToken(TokenFalse)
		}
	case string:
		b.zBuf = strconv.AppendQuote(b.zBuf, v)
		b.setLastToken(TokenString)
	case float32:
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			b.zBuf = append(b.zBuf, "0"...)
		} else {
			b.zBuf = strconv.AppendFloat(b.zBuf, float64(v), 'f', -1, 32)
		}
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			b.zBuf = append(b.zBuf, "0"...)
		} else {
			b.zBuf = strconv.AppendFloat(b.zBuf, v, 'f', -1, 64)
		}
	case int:
		b.zBuf = strconv.AppendInt(b.zBuf, int64(v), 10)
	case int8:
		b.zBuf = strconv.AppendInt(b.zBuf, int64(v), 10)
	case int16:
		b.zBuf = strconv.AppendInt(b.zBuf, int64(v), 10)
	case int32:
		b.zBuf = strconv.AppendInt(b.zBuf, int64(v), 10)
	case int64:
		b.zBuf = strconv.AppendInt(b.zBuf, v, 10)
	case uint:
		b.zBuf = strconv.AppendUint(b.zBuf, uint64(v), 10)
	case uint8:
		b.zBuf = strconv.AppendUint(b.zBuf, uint64(v), 10)
	case uint16:
		b.zBuf = strconv.AppendUint(b.zBuf, uint64(v), 10)
	case uint32:
		b.zBuf = strconv.AppendUint(b.zBuf, uint64(v), 10)
	case uint64:
		b.zBuf = strconv.AppendUint(b.zBuf, v, 10)

	default: // fallback to std encoding/json
		err := stdjson.NewEncoder(b.writer()).Encode(value)
		if err != nil {
			b.addErrorf(err.Error())
		}
		// trim last newline
		if len(b.zBuf) > 0 && b.zBuf[len(b.zBuf)-1] == '\n' {
			b.zBuf = b.zBuf[:len(b.zBuf)-1]
		}
	}
}

func (b *Builder) writeKey(key RawToken) {
	if b.level <= 0 {
		if key.Type() != 0 {
			b.addErrorf("unexpected key(%s) at root", key)
		}
		return
	}
	top := b.stack[len(b.stack)-1]
	switch top {
	case TokenArrayOpen:
		if key.Type() != 0 {
			b.addErrorf("unexpected key(%s) in array", key)
			return
		}
	case TokenObjectOpen:
		if key.Type() == 0 {
			b.addErrorf("missing key in object")
		}
		b.write(key.Raw())
		b.writeByte(':')
		if b.indent != "" {
			b.writeByte(' ')
		}
	default:
		panic("unexpected stack")
	}
}

func (b *Builder) writeKeyString(key string) {
	if b.level <= 0 {
		if key != "" {
			b.addErrorf("unexpected key(%s) at root", key)
		}
		return
	}
	top := b.stack[len(b.stack)-1]
	switch top {
	case TokenArrayOpen:
		if key != "" {
			b.addErrorf("unexpected key(%s) in array", key)
			return
		}
	case TokenObjectOpen:
		if key == "" {
			b.addErrorf("missing key in object")
		}
		b.zBuf = strconv.AppendQuote(b.zBuf, key)
		b.writeByte(':')
		if b.indent != "" {
			b.writeByte(' ')
		}
	default:
		panic("unexpected stack")
	}
}

func (b *Builder) WriteNewline(next TokenType) {
	b.WriteComma(next)
	if b.prefix == "" && b.indent == "" {
		return
	}
	if b.lastNewline {
		return
	}
	if b.lastTok != 0 {
		b.writeByte('\n')
		b.lastNewline = true
	}
}

func (b *Builder) WriteIndent() {
	if b.lastIndent {
		return
	}
	b.lastIndent = true

	b.writeString(b.prefix)
	for range b.level {
		b.writeString(b.indent)
	}
}

func (b *Builder) write(p []byte) {
	b.zBuf = append(b.zBuf, p...)
}

func (b *Builder) writeByte(c byte) {
	b.zBuf = append(b.zBuf, c)
}

func (b *Builder) writeString(s string) {
	b.zBuf = append(b.zBuf, s...)
}

func (b *Builder) writer() *zBuf {
	return &b.zBuf
}

func (b *Builder) setLastToken(tok TokenType) {
	b.lastTok = tok
	b.lastIndent = false
	b.lastNewline = false
}

func (b *Builder) addErrorf(msg string, args ...any) {
	if b.err == nil {
		b.err = fmt.Errorf(msg, args...)
	}
}

// ShouldAddComma returns true if a comma should be added before the next token.
func ShouldAddComma(lastToken, nextToken TokenType) bool {
	switch lastToken {
	case 0, TokenArrayOpen, TokenObjectOpen, TokenComma, TokenColon:
		return false
	default:
		switch nextToken {
		case TokenArrayClose, TokenObjectClose:
			return false
		default:
			return true
		}
	}
}

type zBuf []byte

func (b *zBuf) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}
