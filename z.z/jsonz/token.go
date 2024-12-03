package jsonz

import (
	"fmt"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
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

// GetNumber returns the number value of the token.
func (r RawToken) GetNumber() (float64, error) {
	if r.typ != TokenNumber {
		return 0, fmt.Errorf("invalid number token")
	}
	switch {
	case len(r.raw) == 1 && r.raw[0] == '0':
		return 0, nil
	case len(r.raw) > 1 && r.raw[0] == '0':
		if r.raw[1] == '.' {
			return strconv.ParseFloat(string(r.raw), 64)
		} else {
			return 0, fmt.Errorf("number cannot have leading zero")
		}
	default:
		return strconv.ParseFloat(string(r.raw), 64)
	}
}

// GetBool returns the boolean value of the token.
func (r RawToken) GetBool() (bool, error) {
	switch r.typ {
	case TokenTrue:
		return true, nil
	case TokenFalse:
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean token")
	}
}

// GetString returns the unquoted string value of the token.
// https://datatracker.ietf.org/doc/html/rfc8259#section-7
func (r RawToken) GetString() (string, error) {
	if r.typ != TokenString {
		return "", fmt.Errorf("invalid string token")
	}

	raw := r.raw
	if len(raw) < 2 {
		return "", fmt.Errorf("invalid string token")
	}
	N := len(raw)
	if raw[0] != '"' || raw[N-1] != '"' {
		return "", fmt.Errorf("invalid string token")
	}
	i := 1
	for ; i < N-1; i++ {
		c := raw[i]
		switch c {
		case '\\':
			goto slow
		case '\b', '\f', '\n', '\r', '\t':
			return "", fmt.Errorf("invalid string token")
		}
		if c >= utf8.RuneSelf {
			goto slow // utf-8
		}
	}
	return string(raw[1 : N-1]), nil

slow:
	s := make([]byte, 0, N-2)
	copy(s, raw[1:i])
	N = N - 1 // new length
	for i < N {
		c := raw[i]
		switch c {
		case '\\':
			break
		case '\b', '\f', '\n', '\r', '\t':
			return "", fmt.Errorf("invalid string token")
		default:
			// ascii
			if c < utf8.RuneSelf {
				s = append(s, c)
				i++
				continue
			}
			// utf-8
			r, size := utf8.DecodeRune(raw[i:])
			if r == utf8.RuneError {
				return "", fmt.Errorf("invalid string token")
			}
			s = append(s, raw[i:i+size]...)
			i += size
			continue
		}

		i++
		if i >= N {
			return "", fmt.Errorf("invalid string token")
		}
		switch raw[i] {
		case '"', '\\', '/':
			s = append(s, raw[i])
			i++
		case 'b':
			s = append(s, '\b')
			i++
		case 'f':
			s = append(s, '\f')
			i++
		case 'n':
			s = append(s, '\n')
			i++
		case 'r':
			s = append(s, '\r')
			i++
		case 't':
			s = append(s, '\t')
			i++
		case 'u':
			r, n := decodeHexRune(raw[i-1:])
			if n == 0 {
				return "", fmt.Errorf("invalid string token")
			}
			s = utf8.AppendRune(s, r)
			i += n - 1
		default:
			return "", fmt.Errorf("invalid string token")
		}
	}
	return string(s), nil
}

// decode \uXXXX to rune, return the rune and the number of bytes consumed (0 if error)
func decodeHexRune(s []byte) (rune, int) {
	x := decodeHex(s)
	if x == -1 {
		return utf8.RuneError, 0
	}
	if !utf16.IsSurrogate(x) {
		return x, 6
	}
	xx := decodeHex(s[x+6:])
	if xx == -1 {
		return utf8.RuneError, 0
	}
	r := utf16.DecodeRune(x, xx)
	if r == utf8.RuneError {
		return utf8.RuneError, 0
	}
	return r, 12
}

// decode \uXXXX to rune, return -1 if error (utf8.RuneError is still a valid code point)
func decodeHex(s []byte) (x rune) {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		x = x*16 + rune(c)
	}
	return
}

// GetValue returns the value of the token as an any.
func (r RawToken) GetValue() (any, error) {
	switch r.typ {
	case TokenNull:
		return nil, nil
	case TokenTrue:
		return true, nil
	case TokenFalse:
		return false, nil
	case TokenNumber:
		return r.GetNumber()
	case TokenString:
		return r.GetString()
	case TokenObjectStart, TokenObjectEnd, TokenArrayStart, TokenArrayEnd, TokenComma, TokenColon:
		return r.typ, nil
	}
	return nil, fmt.Errorf("invalid token type: %v", r.typ)
}
