package draft

import (
	"fmt"
	"strings"
)

const rNull = "null"
const rTrue = "true"
const rFalse = "false"

func NextToken(s string) (token string, remain string, err error) {
	s = TrimPrefixSpace(s)
	if s == "" {
		return "", "", nil
	}
	switch s[0] {
	case '{', '}', '[', ']', ',', ':':
		return s[:1], s[1:], nil
	case 'n':
		if len(s) < len(rNull) || s[:len(rNull)] != rNull {
			return "", s, newTokenError(s)
		}
		return rNull, s[len(rNull):], nil
	case 'f':
		if len(s) < len(rFalse) || s[:len(rFalse)] != rFalse {
			return "", s, newTokenError(s)
		}
		return rFalse, s[len(rFalse):], nil
	case 't':
		if len(s) < len(rTrue) || s[:len(rTrue)] != rTrue {
			return "", s, newTokenError(s)
		}
		return rTrue, s[len(rTrue):], nil
	case '"':
		return nextTokenString(s)
	default:
		return nextTokenNumber(s)
	}
}

func nextTokenNumber(s string) (token string, remain string, err error) {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' || c == '+' || c == '-' {
			continue
		}
		if c == '.' || c == 'e' || c == 'E' {
			continue
		}
		token, remain = s[:i], s[i:]
		break
	}
	if token == "" {
		return "", s, newTokenError(s)
	}
	return token, remain, nil
}

func nextTokenString(s string) (token string, remain string, err error) {
	s0, s := s, s[1:]
	for {
		n := strings.IndexByte(s, '"')
		if n < 0 {
			return "", s0, fmt.Errorf("missing close quote at %q", snipStr(s0, 16))
		}
		if s0[n-1] != '\\' {
			return s0[:n+2], s0[n+2:], nil
		}
		i := n - 2
		for ; i >= 0 && s[i] == '\\'; i-- {
		}
		if (n-i)%2 == 0 {
			return s0[:n+2], s0[n+2:], nil
		}
		s = s[n+1:]
	}
}

func TrimPrefixSpace(s string) string {
	if len(s) == 0 || s[0] > 0x20 {
		return s
	}
	if !IsSpace(s[0]) {
		return s
	}
	for i := 1; i < len(s); i++ {
		if !IsSpace(s[i]) {
			return s[i:]
		}
	}
	return ""
}

func IsSpace(c byte) bool {
	return c == 0x09 || c == 0x0A || c == 0x0D || c == 0x20
}

func snipStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// newTokenError returns an error for an invalid token. Input s should not be starts with space.
func newTokenError(s string) error {
	return fmt.Errorf("invalid token at %q", snipStr(s, 16))
}
