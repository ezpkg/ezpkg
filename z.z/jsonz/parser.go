package jsonz

import (
	"fmt"
	"iter"
)

func Parse(in []byte) iter.Seq2[Item, error] {
	return func(yield func(Item, error) bool) {
		defer func() {
			r := recover()
			if r != nil {
				yield(Item{}, fmt.Errorf("%v", r))
			}
		}()

		var tok, next RawToken
		var last *PathItem
		remain := in
		path := make([]PathItem, 0, 16)
		advance := func() {
			var err error
			tok = next
			next, remain, err = NextToken(remain)
			must(err)
		}
		push := func() {
			path = append(path, PathItem{tok: tok})
			last = &path[len(path)-1]
		}
		pop := func() {
			path = path[:len(path)-1]
			if len(path) > 0 {
				last = &path[len(path)-1]
			}
		}
		advance()
		advance()

	value:
		switch {
		case tok.typ == TokenArrayStart:
			push()
			advance()
			switch {
			case tok.typ == TokenArrayEnd:
				goto close
			default:
				goto value
			}

		case tok.typ == TokenObjectStart:
			push()
			advance()
			switch {
			case tok.typ == TokenObjectEnd:
				goto close
			default:
				goto key_value
			}

		case tok.IsValue():
			yield(Item{Path: path, RawToken: tok}, nil)
			switch {
			case last == nil && next.typ == 0:
				return // ✅ done
			case last == nil && next.typ != 0:
				panicf("parser: unexpected token(%s)", tok.typ)
			default:
				advance()
				goto close
			}

		default:
			panicf("parser: invalid token(%s)", tok.typ)
		}

	key_value:
		switch {
		case tok.typ == TokenString:
			last.key = tok
			advance()
			switch {
			case tok.typ == TokenColon:
				advance()
				goto value
			default:
				panicf("parser: unexpected token(%s)", tok.typ)
			}
		default:
			panicf("parser: invalid object start(%s)", tok.typ)
		}

	close:
		switch {
		case tok.typ == TokenArrayEnd:
			if last == nil || last.tok.typ != TokenArrayStart {
				panicf("parser: invalid array end(%s)", tok.typ)
			}
			pop()

		case tok.typ == TokenObjectEnd:
			if last == nil || last.tok.typ != TokenObjectStart {
				panicf("parser: invalid object end(%s)", tok.typ)
			}
			pop()

		case tok.typ == TokenComma:
			last.idx++
			last.key = RawToken{}
			advance()
			switch {
			case last.tok.typ == TokenArrayStart:
				goto value
			case last.tok.typ == TokenObjectStart:
				goto key_value
			default:
				panicf("parser: invalid comma(%s)", tok.typ)
			}
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
