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
	buf    []byte
	indent string
	prefix string

	lastTok TokenType
	level   int
	stack   []TokenType // array or object
	err     error
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
	return b.buf, b.err
}

// AddRaw adds a key and token to the builder. It will add a comma if needed.
func (b *Builder) AddRaw(key, token RawToken) {
	switch {
	case token.IsOpen():
		if ShouldAddComma(b.lastTok, token.Type()) {
			b.writeByte(',')
		}
		b.writeIndent()
		b.writeKey(key)
		b.writeByte(byte(token.Type()))
		b.lastTok = token.Type()
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
		b.writeIndent()
		b.writeByte(byte(token.Type()))
		b.lastTok = token.Type()

	case token.IsValue():
		if ShouldAddComma(b.lastTok, token.Type()) {
			b.writeByte(',')
		}
		b.writeIndent()
		b.writeKey(key)
		b.write(token.Raw())
		b.lastTok = token.Type()

	default:
		b.addErrorf("unexpected token(%s)", token)
	}
}

// AddRaw adds a key and token to the builder. It will add a comma if needed.
func (b *Builder) AddToken(key string, token RawToken) {
	switch {
	case token.IsOpen():
		if ShouldAddComma(b.lastTok, token.Type()) {
			b.writeByte(',')
		}
		b.writeIndent()
		b.writeKeyString(key)
		b.writeByte(byte(token.Type()))
		b.lastTok = token.Type()
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
		b.writeIndent()
		b.writeByte(byte(token.Type()))
		b.lastTok = token.Type()

	case token.IsValue():
		if ShouldAddComma(b.lastTok, token.Type()) {
			b.writeByte(',')
		}
		b.writeIndent()
		b.writeKeyString(key)
		b.write(token.Raw())
		b.lastTok = token.Type()

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
	}

	switch {
	case tokType.IsOpen():
		if ShouldAddComma(b.lastTok, tokType) {
			b.writeByte(',')
		}
		b.writeIndent()
		b.writeKeyString(key)
		b.writeByte(byte(tokType))
		b.lastTok = tokType
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
		b.writeIndent()
		b.writeByte(byte(tokType))
		b.lastTok = tokType

	default:
		if ShouldAddComma(b.lastTok, tokType) {
			b.writeByte(',')
		}
		b.writeIndent()
		b.writeKeyString(key)

		if rawTok.typ != 0 {
			if !rawTok.IsValue() {
				b.addErrorf("unexpected token(%s) as value", rawTok)
				return
			}
			b.write(rawTok.raw)
			b.lastTok = tokType
		} else {
			b.writeValue(value)
		}
	}
}

func (b *Builder) writeValue(value any) {
	b.lastTok = TokenNumber // default to number, can be overridden
	switch v := value.(type) {
	case bool:
		if v {
			b.buf = append(b.buf, "true"...)
			b.lastTok = TokenTrue
		} else {
			b.buf = append(b.buf, "false"...)
			b.lastTok = TokenFalse
		}
	case string:
		b.buf = strconv.AppendQuote(b.buf, v)
		b.lastTok = TokenString
	case float32:
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			b.buf = append(b.buf, "0"...)
		} else {
			b.buf = strconv.AppendFloat(b.buf, float64(v), 'f', -1, 32)
		}
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			b.buf = append(b.buf, "0"...)
		} else {
			b.buf = strconv.AppendFloat(b.buf, v, 'f', -1, 64)
		}
	case int:
		b.buf = strconv.AppendInt(b.buf, int64(v), 10)
	case int8:
		b.buf = strconv.AppendInt(b.buf, int64(v), 10)
	case int16:
		b.buf = strconv.AppendInt(b.buf, int64(v), 10)
	case int32:
		b.buf = strconv.AppendInt(b.buf, int64(v), 10)
	case int64:
		b.buf = strconv.AppendInt(b.buf, v, 10)
	case uint:
		b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
	case uint8:
		b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
	case uint16:
		b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
	case uint32:
		b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
	case uint64:
		b.buf = strconv.AppendUint(b.buf, v, 10)

	default: // fallback to std encoding/json
		err := stdjson.NewEncoder(b.writer()).Encode(value)
		if err != nil {
			b.addErrorf(err.Error())
		}
		// trim last newline
		if len(b.buf) > 0 && b.buf[len(b.buf)-1] == '\n' {
			b.buf = b.buf[:len(b.buf)-1]
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
		b.buf = strconv.AppendQuote(b.buf, key)
		b.writeByte(':')
		if b.indent != "" {
			b.writeByte(' ')
		}
	default:
		panic("unexpected stack")
	}
}

func (b *Builder) writeIndent() {
	if b.prefix == "" && b.indent == "" {
		return
	}
	if b.lastTok != 0 {
		b.writeByte('\n')
	}
	b.writeString(b.prefix)
	for range b.level {
		b.writeString(b.indent)
	}
}

func (b *Builder) write(p []byte) {
	b.buf = append(b.buf, p...)
}

func (b *Builder) writeByte(c byte) {
	b.buf = append(b.buf, c)
}

func (b *Builder) writeString(s string) {
	b.buf = append(b.buf, s...)
}

func (b *Builder) writer() *zBuffer {
	return (*zBuffer)(&b.buf)
}

func (b *Builder) addErrorf(msg string, args ...any) {
	if b.err == nil {
		b.err = fmt.Errorf(msg, args...)
	}
}

// ShouldAddComma returns true if a comma should be added before the next token.
func ShouldAddComma(lastToken, nextToken TokenType) bool {
	switch lastToken {
	case 0, TokenArrayOpen, TokenObjectOpen:
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

type zBuffer []byte

func (b *zBuffer) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}
