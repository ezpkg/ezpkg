package jsoniter

import (
	"bytes"
	"fmt"
	"iter"
)

func NextToken(in []byte) (token RawToken, remain []byte, err error) {
	in = skipSpace(in)
	if len(in) == 0 {
		return RawToken{}, nil, nil
	}
	switch in[0] {
	case '{', '}', '[', ']', ',', ':':
		typ := TokenType(in[0])
		return RawToken{typ: typ, raw: in[:1]}, in[1:], nil
	case 'n':
		return nextTokenConst(in, rNull)
	case 'f':
		return nextTokenConst(in, rFalse)
	case 't':
		return nextTokenConst(in, rTrue)
	case '"':
		return nextTokenString(in)
	default:
		return nextTokenNumber(in)
	}
}

func Scan(in []byte) iter.Seq2[RawToken, error] {
	return func(yield func(token RawToken, err error) bool) {
		remain := in
		for {
			token, rm, err := NextToken(remain)
			remain = rm
			if err != nil {
				yield(RawToken{}, err)
				return
			}
			if !yield(token, nil) {
				return
			}
			if len(remain) == 0 {
				return
			}
		}
	}
}

func nextTokenConst(in []byte, expect []byte) (token RawToken, remain []byte, err error) {
	if len(in) < len(expect) || !bytes.Equal(in[:len(expect)], expect) {
		return RawToken{}, in, newTokenError(string(expect), in)
	}
	typ := TokenType(expect[0])
	return RawToken{typ: typ, raw: expect}, in[len(expect):], nil
}

func nextTokenNumber(in []byte) (token RawToken, remain []byte, err error) {
	i := 0
	for ; i < len(in); i++ {
		c := in[i]
		if c >= '0' && c <= '9' {
			continue
		}
		switch c {
		// https://datatracker.ietf.org/doc/html/rfc8259#section-6
		case '+', '-', '.', 'e', 'E':
			continue
		}
		break
	}
	token, remain = RawToken{typ: TokenNumber, raw: in[:i]}, in[i:]
	if i == 0 {
		return RawToken{}, in, newTokenError("number", in)
	}
	return token, remain, nil
}

func nextTokenString(in []byte) (token RawToken, remain []byte, err error) {
	for i := 1; i < len(in); i++ {
		c := in[i]
		switch c {
		case '"':
			return RawToken{typ: TokenString, raw: in[:i+1]}, in[i+1:], nil
		case '\\':
			if i+1 >= len(in) {
				return RawToken{}, in, newTokenError("string", in)
			}
			switch in[i+1] {
			// https://datatracker.ietf.org/doc/html/rfc8259#section-7
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				i++
			case 'u':
				if i+5 >= len(in) {
					return RawToken{}, in, newTokenError("string", in)
				}
				i += 5
			default:
				return RawToken{}, in, newTokenError("string", in)
			}
		}
	}
	return RawToken{}, in, newTokenError("string", in)
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
	return fmt.Errorf("expect %s, got %s", name, snip(in, 16))
}
