package jsonz

import (
	"bytes"
	"fmt"
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
	bytes.Buffer
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
	return b.Buffer.Bytes(), b.err
}

// AddRaw adds a key and token to the builder. It will add a comma if needed.
func (b *Builder) AddRaw(key, token RawToken) {
	switch {
	case token.IsOpen():
		if ShouldAddComma(b.lastTok, token.Type()) {
			b.WriteByte(',')
		}
		b.writeIndent()
		b.writeKey(key)
		b.WriteByte(byte(token.Type()))
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
		b.WriteByte(byte(token.Type()))
		b.lastTok = token.Type()

	case token.IsValue():
		if ShouldAddComma(b.lastTok, token.Type()) {
			b.WriteByte(',')
		}
		b.writeIndent()
		b.writeKey(key)
		b.Write(token.Raw())
		b.lastTok = token.Type()
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
		b.Write(key.Raw())
		b.WriteByte(':')
		if b.indent != "" {
			b.WriteByte(' ')
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
		b.WriteByte('\n')
	}
	b.WriteString(b.prefix)
	for range b.level {
		b.WriteString(b.indent)
	}
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
