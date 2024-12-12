package jsoniter

import (
	"bytes"
	"math"
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
	TokenObjectOpen  TokenType = '{'
	TokenObjectClose TokenType = '}'
	TokenArrayOpen   TokenType = '['
	TokenArrayClose  TokenType = ']'
	TokenComma       TokenType = ','
	TokenColon       TokenType = ':'
)

var (
	tokenNull        = RawToken{typ: TokenNull, raw: []byte(`null`)}
	tokenTrue        = RawToken{typ: TokenTrue, raw: []byte(`true`)}
	tokenFalse       = RawToken{typ: TokenFalse, raw: []byte(`false`)}
	tokenObjectOpen  = RawToken{typ: TokenObjectOpen, raw: []byte(`{`)}
	tokenObjectClose = RawToken{typ: TokenObjectClose, raw: []byte(`}`)}
	tokenArrayOpen   = RawToken{typ: TokenArrayOpen, raw: []byte(`[`)}
	tokenArrayClose  = RawToken{typ: TokenArrayClose, raw: []byte(`]`)}
	tokenComma       = RawToken{typ: TokenComma, raw: []byte(`,`)}
	tokenColon       = RawToken{typ: TokenColon, raw: []byte(`:`)}
	tokenString      = RawToken{typ: TokenString, raw: []byte(`""`)}
	tokenNumber      = RawToken{typ: TokenNumber, raw: []byte(`0`)}
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

// New returns a new raw token of the type. For number and string, use NumberToken() and StringToken().
func (t TokenType) New() RawToken {
	switch t {
	case TokenNull:
		return tokenNull
	case TokenTrue:
		return tokenTrue
	case TokenFalse:
		return tokenFalse
	case TokenNumber:
		return tokenNumber
	case TokenString:
		return tokenString
	case TokenObjectOpen:
		return tokenObjectOpen
	case TokenObjectClose:
		return tokenObjectClose
	case TokenArrayOpen:
		return tokenArrayOpen
	case TokenArrayClose:
		return tokenArrayClose
	case TokenComma:
		return tokenComma
	case TokenColon:
		return tokenColon
	default:
		return RawToken{}
	}
}

// IsOpen returns true if the token is an open token '[' or '{'.
func (t TokenType) IsOpen() bool {
	return t == TokenArrayOpen || t == TokenObjectOpen
}

// IsClose returns true if the token is a close token ']' or '}'.
func (t TokenType) IsClose() bool {
	return t == TokenArrayClose || t == TokenObjectClose
}

// IsValue returns true if the token is a value: null, boolean, number, string.
func (t TokenType) IsValue() bool {
	switch t {
	case TokenNull, TokenTrue, TokenFalse, TokenNumber, TokenString:
		return true
	default:
		return false
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

// NewRawToken returns a new raw token from the raw bytes.
func NewRawToken(raw []byte) (RawToken, error) {
	token, remain, err := newRawToken(raw)
	switch {
	case err != nil:
		return RawToken{}, err
	case token.typ == 0:
		return RawToken{}, ErrTokenEmpty
	case len(remain) != 0:
		return RawToken{}, ErrTokenInvalid
	}
	return token, err
}

func newRawToken(raw []byte) (RawToken, []byte, error) {
	token, remain, err := NextToken(raw)
	if err != nil {
		return RawToken{}, remain, err
	}
	switch token.typ {
	case TokenNumber:
		_, err = token.GetNumber() // validate number
	case TokenString:
		_, err = token.GetString() // validate string
	}
	return token, remain, err
}

// MustRawToken returns a new raw token from the raw bytes. Panic if error.
func MustRawToken(raw []byte) RawToken {
	token, err := NewRawToken(raw)
	must(err)
	return token
}

// StringToken returns a string token.
func StringToken(s string) RawToken {
	return RawToken{typ: TokenString, raw: quoteString(s)}
}

// NumberToken returns a number token. For NaN and Inf, fallback to 0.
func NumberToken(f float64) RawToken {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return RawToken{typ: TokenNumber, raw: []byte("0")}
	}
	return RawToken{typ: TokenNumber, raw: []byte(strconv.FormatFloat(f, 'g', -1, 64))}
}

// IntToken returns a number token.
func IntToken(n int) RawToken {
	return RawToken{typ: TokenNumber, raw: []byte(strconv.Itoa(n))}
}

// BoolToken returns a boolean token.
func BoolToken(b bool) RawToken {
	if b {
		return tokenTrue
	}
	return tokenFalse
}

// Type returns the type of the token.
func (r RawToken) Type() TokenType {
	return r.typ
}

// Bytes returns the raw bytes value of the token.
func (r RawToken) Bytes() []byte {
	return r.raw
}

// Raw returns the raw bytes value of the token.
func (r RawToken) Raw() []byte {
	return r.raw
}

// String returns the raw string value of the token. Use ToString() for unquoted strings.
func (r RawToken) String() string {
	return string(r.raw)
}

// IsZero returns true if the token is zero.
func (r RawToken) IsZero() bool {
	return r.typ == 0
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

// IsOpen returns true if the token is an open token '[' or '{'.
func (r RawToken) IsOpen() bool {
	return r.typ == TokenArrayOpen || r.typ == TokenObjectOpen
}

// IsClose returns true if the token is a close token ']' or '}'.
func (r RawToken) IsClose() bool {
	return r.typ == TokenArrayClose || r.typ == TokenObjectClose
}

// GetNumber returns the number value of the token.
func (r RawToken) GetNumber() (float64, error) {
	if r.typ != TokenNumber {
		return 0, ErrTokenNumber
	}
	switch {
	case len(r.raw) == 1 && r.raw[0] == '0':
		return 0, nil
	case len(r.raw) > 1 && r.raw[0] == '0':
		if r.raw[1] == '.' {
			return strconv.ParseFloat(string(r.raw), 64)
		} else {
			// number cannot have leading zero
			return 0, ErrTokenNumber
		}
	default:
		f, err := strconv.ParseFloat(string(r.raw), 64)
		switch {
		case err != nil:
			return 0, ErrTokenNumber
		case math.IsNaN(f) || math.IsInf(f, 0):
			return 0, ErrTokenNumber
		}
		return f, nil
	}
}

// GetInt returns the integer value of the token.
func (r RawToken) GetInt() (int, error) {
	if r.typ != TokenNumber {
		return 0, ErrTokenNumber
	}
	if bytes.ContainsAny(r.raw, ".eE") {
		return 0, ErrNumberNotInt
	}
	v, err := strconv.ParseInt(string(r.raw), 10, 64)
	if err != nil {
		return 0, ErrTokenNumber
	}
	return int(v), nil
}

// GetBool returns the boolean value of the token.
func (r RawToken) GetBool() (bool, error) {
	switch r.typ {
	case TokenTrue:
		return true, nil
	case TokenFalse:
		return false, nil
	default:
		return false, ErrTokenBool
	}
}

// GetString returns the unquoted string value of the token.
// https://datatracker.ietf.org/doc/html/rfc8259#section-7
func (r RawToken) GetString() (string, error) {
	if r.typ != TokenString {
		return "", ErrTokenString
	}
	b, err := unquote(r.raw)
	return string(b), err
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
	case TokenObjectOpen, TokenObjectClose, TokenArrayOpen, TokenArrayClose, TokenComma, TokenColon:
		return r.typ, nil
	}
	return nil, ErrTokenType
}

// Equal returns true if the token is equal to the other token.
func (r RawToken) Equal(other RawToken) bool {
	return r.typ == other.typ && string(r.raw) == string(other.raw)
}

func unquote(raw []byte) ([]byte, error) {
	N := len(raw)
	if N < 2 {
		return nil, ErrTokenString
	}
	if raw[0] != '"' || raw[N-1] != '"' {
		return nil, ErrTokenString
	}
	i := 1
	for ; i < N-1; i++ {
		c := raw[i]
		switch c {
		case '\\':
			goto slow
		case '\b', '\f', '\n', '\r', '\t':
			return nil, ErrTokenString
		}
		if c >= utf8.RuneSelf {
			goto slow // utf-8
		}
	}
	return raw[1 : N-1], nil

slow:
	b := make([]byte, 0, N-2)
	copy(b, raw[1:i])
	N = N - 1 // new length
	for i < N {
		c := raw[i]
		switch c {
		case '\\':
			break
		case '\b', '\f', '\n', '\r', '\t':
			return nil, ErrTokenString
		default:
			// ascii
			if c < utf8.RuneSelf {
				b = append(b, c)
				i++
				continue
			}
			// utf-8
			r, size := utf8.DecodeRune(raw[i:])
			if r == utf8.RuneError {
				return nil, ErrTokenString
			}
			b = append(b, raw[i:i+size]...)
			i += size
			continue
		}

		i++
		if i >= N {
			return nil, ErrTokenString
		}
		switch raw[i] {
		case '"', '\\', '/':
			b = append(b, raw[i])
			i++
		case 'b':
			b = append(b, '\b')
			i++
		case 'f':
			b = append(b, '\f')
			i++
		case 'n':
			b = append(b, '\n')
			i++
		case 'r':
			b = append(b, '\r')
			i++
		case 't':
			b = append(b, '\t')
			i++
		case 'u':
			utfRune, n := decodeHexRune(raw[i-1:])
			if n == 0 {
				return nil, ErrTokenString
			}
			b = utf8.AppendRune(b, utfRune)
			i += n - 1
		default:
			return nil, ErrTokenString
		}
	}
	return b, nil
}

func canSimplyUnquote(raw []byte) bool {
	N := len(raw)
	if N <= 2 || raw[0] != '"' || raw[N-1] != '"' {
		return false
	}
	for i := 1; i < N-1; i++ {
		c := raw[i]
		if c > 0x20 && c <= 0x7E && c != '"' && c != '\\' {
			// printable ascii
		} else {
			return false
		}
	}
	return true
}

func needQuote(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c > 0x7E {
			return true
		}
		switch c {
		case '"', '\\':
			return true
		}
	}
	return false
}

func quoteString(s string) []byte {
	var b []byte
	return strconv.AppendQuote(b, s)
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

func intOrStr(x any) any {
	switch x := x.(type) {
	case int:
		return x
	case string:
		return x
	case int8:
		return int(x)
	case int16:
		return int(x)
	case int32:
		return int(x)
	case int64:
		return int(x)
	case uint:
		return int(x)
	case uint8:
		return int(x)
	case uint16:
		return int(x)
	case uint32:
		return int(x)
	case uint64:
		return int(x)
	default:
		return x
	}
}
