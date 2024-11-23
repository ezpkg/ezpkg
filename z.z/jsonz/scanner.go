package draft

import (
	"bytes"
	"fmt"
)

var rNull = []byte("null")
var rTrue = []byte("true")
var rFalse = []byte("false")

type RawToken []byte

func NextToken(in []byte) (token RawToken, remain []byte, err error) {
	in = skipSpace(in)
	if len(in) == 0 {
		return nil, nil, nil
	}
	switch in[0] {
	case '{', '}', '[', ']', ',', ':':
		return in[:1], in[1:], nil
	case 'n':
		return nextToken(in, rNull)
	case 'f':
		return nextToken(in, rFalse)
	case 't':
		return nextToken(in, rTrue)
	case '"':
		return nextTokenString(in)
	default:
		return nextTokenNumber(in)
	}
}

func nextToken(in []byte, expect []byte) (token RawToken, remain []byte, err error) {
	if len(in) < len(expect) || !bytes.Equal(in[:len(expect)], expect) {
		return nil, in, newTokenError(string(expect), in)
	}
	return expect, in[len(expect):], nil
}

func nextTokenNumber(in []byte) (token RawToken, remain []byte, err error) {
	for i := 0; i < len(in); i++ {
		c := in[i]
		if c >= '0' && c <= '9' {
			continue
		}
		switch c {
		// https://datatracker.ietf.org/doc/html/rfc8259#section-6
		case '+', '-', '.', 'e', 'E':
			continue
		}
		token, remain = in[:i], in[i:]
		break
	}
	if len(token) == 0 {
		return nil, in, newTokenError("number", in)
	}
	return token, remain, nil
}

func nextTokenString(in []byte) (token RawToken, remain []byte, err error) {
	for i := 1; i < len(in); i++ {
		c := in[i]
		switch c {
		case '"':
			return in[:i+1], in[i+1:], nil
		case '\\':
			if i+1 >= len(in) {
				return nil, in, newTokenError("string", in)
			}
			switch in[i+1] {
			// https://datatracker.ietf.org/doc/html/rfc8259#section-7
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				i++
			case 'u':
				if i+5 >= len(in) {
					return nil, in, newTokenError("string", in)
				}
				i += 5
			default:
				return nil, in, newTokenError("string", in)
			}
		}
	}
	return nil, in, newTokenError("string", in)
}

func skipSpace(in []byte) []byte {
	for i := 0; i < len(in); i++ {
		c := in[i]
		switch c {
		// https://datatracker.ietf.org/doc/html/rfc8259#section-2
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return in[i:]
		}
	}
	return nil
}

func snip(in []byte, n int) []byte {
	if len(in) <= n {
		return in
	}
	return in[:n]
}

// newTokenError returns an error for an invalid token. Input should not be starts with space.
func newTokenError(name string, in []byte) error {
	return fmt.Errorf("expect %v, got %v", name, snip(in, 16))
}
