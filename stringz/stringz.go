package stringz // import "ezpkg.io/stringz"
import (
	"strings"
)

func FirstNParts(s string, n int, sep string) string {
	idx := FirstNPartsIdx(s, n, sep)
	return s[:idx]
}
func FirstNPartsX(s string, n int, sep string) (parts, remain string) {
	idx := FirstNPartsIdx(s, n, sep)
	return s[:idx], strings.TrimPrefix(s[idx:], sep)
}
func FirstNPartsIdx(s string, n int, sep string) int {
	if n == 0 {
		return 0
	}
	for i := 0; i < len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			n--
			if n <= 0 {
				return i
			}
		}
	}
	return len(s)
}

func LastNParts(s string, n int, sep string) string {
	idx := LastNPartsIdx(s, n, sep)
	return s[idx:]
}
func LastNPartsX(s string, n int, sep string) (parts, remain string) {
	idx := LastNPartsIdx(s, n, sep)
	return s[idx:], strings.TrimSuffix(s[:idx], sep)
}
func LastNPartsIdx(s string, n int, sep string) int {
	if n == 0 {
		return len(s)
	}
	for i := len(s) - len(sep) - 1; i >= 0; i-- {
		if s[i:i+len(sep)] == sep {
			n--
			if n <= 0 {
				return i + len(sep)
			}
		}
	}
	return 0
}

func FirstPart(s, sep string) string {
	return FirstNParts(s, 1, sep)
}
func FirstPartX(s, sep string) (first, remain string) {
	return FirstNPartsX(s, 1, sep)
}
func FirstPartIdx(s, sep string) int {
	return FirstNPartsIdx(s, 1, sep)
}
func LastPart(s, sep string) string {
	return LastNParts(s, 1, sep)
}
func LastPartX(s, sep string) (last, remain string) {
	return LastNPartsX(s, 1, sep)
}
func LastPartIdx(s, sep string) int {
	return LastNPartsIdx(s, 1, sep)
}
