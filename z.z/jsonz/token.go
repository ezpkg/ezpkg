package jsonz

import (
	"fmt"
	"strconv"
)

// TokenType represents the type of a JSON token.
type TokenType byte

const (
	TokenNull        TokenType = 'n'
	TokenTrue        TokenType = 't'
	TokenFalse       TokenType = 'f'
	TokenNumber      TokenType = '0'
	TokenString      TokenType = '"'
	TokenObjectStart TokenType = '{'
	TokenObjectEnd   TokenType = '}'
	TokenArrayStart  TokenType = '['
	TokenArrayEnd    TokenType = ']'
	TokenComma       TokenType = ','
	TokenColon       TokenType = ':'
)

func (t TokenType) String() string {
	switch t {
	case 0:
		return "EOF"
	case TokenNull:
		return "null"
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"
	case TokenNumber:
		return "number"
	case TokenString:
		return "string"
	default:
		return string(t)
	}
}

// RawToken represents a raw token from the scanner.
type RawToken struct {
	typ TokenType
	raw []byte
}

var (
	rNull  = []byte("null")
	rTrue  = []byte("true")
	rFalse = []byte("false")
)

// Type returns the type of the token.
func (r RawToken) Type() TokenType {
	return r.typ
}

// Bytes returns the raw bytes value of the token.
func (r RawToken) Bytes() []byte {
	return r.raw
}

// String returns the raw string value of the token. Use ToString() for unquoted strings.
func (r RawToken) String() string {
	return string(r.raw)
}

// IsValue returns true if the token is a value.
func (r RawToken) IsValue() bool {
	switch r.typ {
	case TokenNull, TokenTrue, TokenFalse, TokenNumber, TokenString:
		return true
	default:
		return false
	}
}

// ParseNumber returns the number value of the token.
func (r RawToken) ParseNumber() (float64, error) {
	switch r.typ {
	case TokenNumber:
		return strconv.ParseFloat(string(r.raw), 64)
	default:
		return 0, fmt.Errorf("invalid number token: %v", r.typ)
	}
}

// ParseBool returns the boolean value of the token.
func (r RawToken) ParseBool() (bool, error) {
	switch r.typ {
	case TokenTrue:
		return true, nil
	case TokenFalse:
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean token: %v", r.typ)
	}
}

// ParseString returns the unquoted string value of the token.
func (r RawToken) ParseString() (string, error) {
	switch r.typ {
	case TokenString:
		return strconv.Unquote(string(r.raw))
	default:
		return "", fmt.Errorf("invalid string token: %v", r.typ)
	}
}

// ParseAny returns the value of the token as an any.
func (r RawToken) ParseAny() (any, error) {
	switch r.typ {
	case TokenNull:
		return nil, nil
	case TokenTrue:
		return true, nil
	case TokenFalse:
		return false, nil
	case TokenNumber:
		return r.ParseNumber()
	case TokenString:
		return r.ParseString()
	case TokenObjectStart, TokenObjectEnd, TokenArrayStart, TokenArrayEnd, TokenComma, TokenColon:
		return r.typ, nil
	}
	return nil, fmt.Errorf("invalid token type: %v", r.typ)
}
