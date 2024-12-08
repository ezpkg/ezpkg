package jsonz

import (
	"bytes"

	"ezpkg.io/bytez"
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
	b := bytez.Buffer{}
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
			b.Println(prefix)
		}
		for range item.Level {
			b.WriteString(indent)
		}
		if item.Key.IsValue() {
			b.Write(item.Key.Raw())
			b.WriteByte(':')
		}
		b.Write(item.Token.Raw())
		lastToken = item.Token.Type()
	}
	return b.Bytes(), nil
}

func ShouldAddComma(lastToken, nextToken TokenType) bool {
	switch lastToken {
	case 0, TokenArrayStart, TokenObjectStart:
		return false
	default:
		switch nextToken {
		case TokenArrayEnd, TokenObjectEnd:
			return false
		default:
			return true
		}
	}
}
