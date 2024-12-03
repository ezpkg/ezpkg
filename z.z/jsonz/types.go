package jsonz

import (
	"fmt"

	"ezpkg.io/fmtz"
)

type Item struct {
	RawToken
	Path Path
}

type Path []PathItem

type PathItem struct {
	idx int
	tok RawToken
	key RawToken
}

func (x Item) String() string {
	return x.RawToken.String()
}
func (x Item) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	switch {
	case f.Flag('+'):
		fz.Printf("%s → %s", x.Path, x.RawToken)

	default:
		fz.Print(x.RawToken)
	}
}
func (x Item) ParseValue() (any, error) {
	return x.RawToken.ParseValue()
}

func (p Path) String() string {
	return fmt.Sprint(p)
}
func (p Path) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	for _, item := range p {
		fz.Printf("%s", item)
	}
}

func (p PathItem) IsArray() bool {
	return p.tok.typ == TokenArrayStart
}
func (p PathItem) IsObject() bool {
	return p.tok.typ == TokenObjectStart
}
func (p PathItem) String() string {
	return fmt.Sprint(p)
}
func (p PathItem) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	switch {
	case p.IsArray():
		fz.Printf("[%d]", p.idx)
	case p.IsObject():
		fz.Printf(".%s", p.key)
	default:
		// root
	}
}
