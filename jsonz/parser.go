package jsonz

import (
	"fmt"
	"iter"
)

func Parse(in []byte) iter.Seq2[Item, error] {
	return func(yield func(Item, error) bool) {
		defer func() {
			if r := recover(); r != nil {
				yield(Item{}, fmt.Errorf("%v", r))
			}
		}()

		var tok, next RawToken
		remain := in
		path := make([]PathItem, 1, 16)
		last := &path[0]
		advance := func() {
			var err error
			tok = next
			next, remain, err = NextToken(remain)
			must(err)
		}
		push := func() {
			path = append(path, PathItem{Token: tok})
			last = &path[len(path)-1]
		}
		pop := func() {
			path = path[:len(path)-1]
			last = &path[len(path)-1]
		}
		yieldValue := func(key RawToken) bool {
			item := Item{
				Level: len(path) - 1, Index: last.Index,
				Key: key, Token: tok, path: path,
			}
			return yield(item, nil)
		}
		advance()
		advance()

	value:
		switch {
		case tok.typ == TokenArrayStart:
			if !yieldValue(last.Key) {
				return
			}
			push()
			advance()
			if tok.typ == TokenArrayEnd {
				goto close
			} else {
				goto value
			}

		case tok.typ == TokenObjectStart:
			if !yieldValue(last.Key) {
				return
			}
			push()
			advance()
			if tok.typ == TokenObjectEnd {
				goto close
			} else {
				goto key_value
			}

		case tok.IsValue():
			if !yieldValue(last.Key) {
				return
			}
			switch {
			case last.Token.typ == 0 && next.typ == 0:
				return // ✅ done
			case last.Token.typ == 0 && next.typ != 0:
				panicf("parser: unexpected token(%s)", tok.typ)
			default:
				advance()
				goto close
			}

		default:
			panicf("parser: expected value, got(%s)", tok.typ)
		}

	key_value:
		switch {
		case tok.typ == TokenString:
			last.Key = tok
			advance()
			if tok.typ == TokenColon {
				advance()
				goto value
			} else {
				panicf("parser: expected colon, got(%s)", tok.typ)
			}
		default:
			panicf("parser: expected key, got(%s)", tok.typ)
		}

	close:
		switch {
		case tok.typ == TokenArrayEnd:
			if last.Token.typ != TokenArrayStart {
				panicf("parser: unexpected array end")
			}
			pop()
			if !yieldValue(RawToken{}) {
				return
			}
			advance()
			if len(path) > 1 {
				goto close
			} else {
				goto end
			}

		case tok.typ == TokenObjectEnd:
			if last.Token.typ != TokenObjectStart {
				panicf("parser: unexpected object end")
			}
			pop()
			if !yieldValue(RawToken{}) {
				return
			}
			advance()
			if len(path) > 1 {
				goto close
			} else {
				goto end
			}

		case tok.typ == TokenComma:
			last.Index++
			last.Key = RawToken{}
			advance()
			switch {
			case last.Token.typ == TokenArrayStart:
				goto value
			case last.Token.typ == TokenObjectStart:
				goto key_value
			default:
				panicf("parser: unexpected comma")
			}

		default:
			panicf("parser: unexpected token(%s)", tok.typ)
		}

	end:
		if tok.typ != 0 {
			panicf("parser: unexpected token(%s)", tok.typ)
		}
	}
}

func panicf(format string, args ...any) {
	panic(fmt.Errorf(format, args...))
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
