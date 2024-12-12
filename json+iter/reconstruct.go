package jsoniter

import (
	"bytes"
	stdjson "encoding/json"
	"iter"
	"math"
	"strconv"

	"ezpkg.io/errorz"
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
	altBuf []byte

	indent string
	prefix string

	stItem
	level int
	stack []stItem // array or object
	err   error

	skipEmptyStructures bool
	useAltBuf           bool
}

// NewBuilder creates a new Builder. It's optional to set the prefix and indent. A zero Builder is valid.
func NewBuilder(prefix, indent string) *Builder {
	return &Builder{
		indent: indent,
		prefix: prefix,
	}
}

// SetSkipEmptyStructures makes the Builder ignores empty array `[]` and empty object `{}`. Default to false.
func (b *Builder) SetSkipEmptyStructures(skip bool) {
	b.skipEmptyStructures = skip
	if skip {
		b.useAltBuf = true
		if b.altBuf == nil {
			b.altBuf = make([]byte, 0, 64)
		}
	} else {
		b.switchBuf()
	}
}

// Bytes returns the bytes of the builder with an error if any.
func (b *Builder) Bytes() ([]byte, error) {
	if b.Len() == 0 {
		return []byte("null"), nil
	}
	return append(b.buf, b.altBuf...), b.err
}

// AddRaw adds a key and token to the builder. It will add a comma if needed.
func (b *Builder) AddRaw(key, token RawToken) {
	b.add(key, token.typ, token.raw, nil)
}

// Add adds a key and value to the builder. It will add a comma if needed.
func (b *Builder) Add(key any, value any) {
	var token RawToken
	var parseNext func() (Item, error, bool)
	var stopParse func()

	switch v := value.(type) {
	case TokenType: // open or close token
		token = v.New()

	case RawToken:
		token = v

	case []byte:
		tok, remain, err := newRawToken(v)
		if tok.typ != 0 && err == nil && len(remain) == 0 {
			token = tok // valid token
		} else {
			parseNext, stopParse = iter.Pull2(Parse(v))
			defer stopParse()
			item, err0, ok := parseNext()
			switch {
			case err0 != nil:
				b.addErrorf("Add: %v", err0)
			case !ok || item.Token.Type() == 0:
				b.addErrorf("Add: empty value")
			default:
				token = item.Token
			}
		}

	case stdjson.RawMessage:
		var err error
		token, err = NewRawToken(v)
		if err != nil {
			b.addErrorf(err.Error())
		}

	default:
		// string | number | boolean | null
	}

	b.add(key, token.typ, token.raw, value)
	if parseNext == nil {
		return
	}
	for item, err0, ok := parseNext(); ok; item, err0, ok = parseNext() {
		if err0 != nil {
			b.addErrorf("Add: %v", err0)
			return
		}
		b.add(item.Key, item.Token.Type(), item.Token.Raw(), nil)
	}
}

func (b *Builder) add(key any, tokType TokenType, raw []byte, value any) {
	switch {
	case tokType.IsOpen():
		snapshot := b.snapshot()
		b.switchAltBuf()
		b.WriteNewline(tokType)
		b.WriteIndent()
		b.writeKey(key)
		b.writeByte(byte(tokType))
		b.push(tokType, snapshot)
		b.setLastToken(tokType)

	case tokType.IsClose():
		if isValidKey(key) {
			b.addErrorf("unexpected key(%s) before close token(%s)", key, tokType)
			return
		}
		snapshot, ok := b.pop()
		if !ok {
			b.addErrorf("unexpected close token(%s)", tokType)
			return
		}
		if b.skipEmptyStructures && b.lastTok.IsOpen() {
			b.restore(snapshot)
		} else {
			b.WriteNewline(tokType)
			b.WriteIndent()
			b.writeByte(byte(tokType))
			b.setLastToken(tokType)
		}

	case tokType == 0 && raw == nil: // value is set
		switch intOrStr(value).(type) {
		case int, float32, float64:
			tokType = TokenNumber
		case string:
			tokType = TokenString
		case bool:
			tokType = TokenTrue
		case nil:
			tokType = TokenNull
		default:
			// fallback to encoding/json
			tokType = TokenString
		}
		fallthrough

	case tokType.IsValue():
		b.switchBuf()
		b.WriteNewline(tokType)
		b.WriteIndent()
		b.writeKey(key)
		if raw != nil {
			b.write(raw)
		} else {
			b.writeValue(value)
		}
		b.setLastToken(tokType)

	default:
		if raw != nil {
			b.addErrorf("unexpected token(%s)", raw)
		} else {
			b.addErrorf("unexpected token(%s)", value)
		}
	}
}

// ShouldAddComma returns true if a comma should be added before the next token.
func (b *Builder) ShouldAddComma(next TokenType) bool {
	return ShouldAddComma(b.lastTok, next)
}

// WriteCommand writes a comma if needed.
func (b *Builder) WriteComma(next TokenType) {
	ok := b.ShouldAddComma(next) || next == 0
	if ok {
		b.writeByte(',')
		b.setLastToken(TokenComma)
	}
	return
}

func (b *Builder) writeValue(value any) {
	b.setLastToken(TokenNumber) // default to number, can be overridden
	switch v := value.(type) {
	case nil:
		b.writeString("null")
		b.setLastToken(TokenNull)
	case bool:
		if v {
			b.buf = append(b.buf, "true"...)
			b.setLastToken(TokenTrue)
		} else {
			b.buf = append(b.buf, "false"...)
			b.setLastToken(TokenFalse)
		}
	case string:
		b.buf = strconv.AppendQuote(b.buf, v)
		b.setLastToken(TokenString)
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

	default:
		// fallback to std encoding/json (TODO: need more work here)
		err := stdjson.NewEncoder(b).Encode(value)
		if err != nil {
			b.addErrorf(err.Error())
		}
		// trim last newline
		if len(b.buf) > 0 && b.buf[len(b.buf)-1] == '\n' {
			b.buf = b.buf[:len(b.buf)-1]
		}
	}
}

func (b *Builder) writeKey(key any) {
	var keyType TokenType
	var raw []byte
	var str string

	switch key := key.(type) {
	case RawToken:
		keyType = key.Type()
		raw = key.Raw()
	case string:
		keyType = TokenString
		str = key
	default:
		b.addErrorf("writeKey: unexpected key type(%T)", key)
		return
	}

	if b.level <= 0 {
		if isValidKey(key) {
			b.addErrorf("writeKey: unexpected key(%s) at root", key)
		}
		return
	}
	top := b.stack[len(b.stack)-1]
	switch top.tok {
	case TokenArrayOpen:
		if isValidKey(key) {
			b.addErrorf("writeKey: unexpected key(%s) in array", key)
			return
		}
	case TokenObjectOpen:
		if keyType == 0 {
			b.addErrorf("writeKey: missing key in object")
		}
		if raw != nil {
			b.write(raw)
		} else {
			b.writeQuote(str)
		}
		b.writeByte(':')
		if b.indent != "" {
			b.writeByte(' ')
		}
	default:
		b.addErrorf("writeKey: unexpected stack")
	}
}

// WriteNewline writes a newline if needed.
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

// WriteIndent writes the indentation if needed.
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

func (b *Builder) Write(p []byte) (n int, err error) {
	if b.useAltBuf {
		b.altBuf = append(b.altBuf, p...)
	} else {
		b.buf = append(b.buf, p...)
	}
	return len(p), nil
}

func (b *Builder) write(p []byte) {
	if b.useAltBuf {
		b.altBuf = append(b.altBuf, p...)
	} else {
		b.buf = append(b.buf, p...)
	}
}

func (b *Builder) writeByte(c byte) {
	if b.useAltBuf {
		b.altBuf = append(b.altBuf, c)
	} else {
		b.buf = append(b.buf, c)
	}
}

func (b *Builder) writeString(s string) {
	if b.useAltBuf {
		b.altBuf = append(b.altBuf, s...)
	} else {
		b.buf = append(b.buf, s...)
	}
}

func (b *Builder) writeQuote(key string) {
	if b.useAltBuf {
		b.altBuf = strconv.AppendQuote(b.altBuf, key)
	} else {
		b.buf = strconv.AppendQuote(b.buf, key)
	}
}

func (b *Builder) switchBuf() {
	b.buf = append(b.buf, b.altBuf...)
	b.altBuf = b.altBuf[:0]
	b.useAltBuf = false
}

func (b *Builder) switchAltBuf() {
	b.useAltBuf = true
}

func (b *Builder) snapshot() stItem {
	snapshot := b.stItem
	snapshot.idx = b.Len()
	return snapshot
}

func (b *Builder) restore(snapshot stItem) {
	altIdx := snapshot.idx - len(b.buf)
	if altIdx >= 0 {
		b.altBuf = b.altBuf[:altIdx]
	}
	b.stItem = snapshot
}

func (b *Builder) push(token TokenType, snapshot stItem) {
	b.level++
	item := snapshot
	item.tok = token
	b.stack = append(b.stack, item)
}

func (b *Builder) pop() (stItem, bool) {
	if len(b.stack) == 0 {
		b.addErrorf("unexpected: pop with zero stack")
		return stItem{}, false
	}
	b.level--
	top := b.stack[len(b.stack)-1]
	b.stack = b.stack[:len(b.stack)-1]
	return top, true
}

func (b *Builder) peek() (stItem, bool) {
	if len(b.stack) == 0 {
		return stItem{}, false
	}
	return b.stack[len(b.stack)-1], true
}

func (b *Builder) setLastToken(tok TokenType) {
	b.lastTok = tok
	b.lastIndent = false
	b.lastNewline = false
}

// Len returns the number of bytes written.
func (b *Builder) Len() int {
	return len(b.buf) + len(b.altBuf)
}

func (b *Builder) addErrorf(msg string, args ...any) {
	if b.err == nil {
		b.err = errorz.Newf(msg, args...)
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

type stItem struct {
	idx int
	tok TokenType

	lastTok     TokenType
	lastIndent  bool
	lastNewline bool
}

func isValidKey(key any) bool {
	switch key := key.(type) {
	case string:
		return key != ""
	case RawToken:
		return key.IsValue()
	default:
		return false
	}
}
