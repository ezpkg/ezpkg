package jsonz

import (
	"fmt"

	"ezpkg.io/fmtz"
)

type Item struct {
	Level int      // level of indentation
	Index int      // index in the parent array or object
	Key   RawToken // optional object "key"
	Token RawToken // [ or { or } or ] or , or value

	path RawPath
}

type RawPath []PathItem

type PathItem struct {
	Index int
	Token RawToken
	Key   RawToken
}

func (x Item) String() string {
	return x.Token.String()
}

func (x Item) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	switch {
	case f.Flag('+'):
		fz.Printf("%s → %s", x.path, x.Token)

	default:
		fz.Print(x.Token)
	}
}

func (x Item) IsArrayValue() bool {
	return !x.Key.IsValue() && x.Token.IsValue()
}

func (x Item) IsObjectValue() bool {
	return x.Key.IsValue() && x.Token.IsValue()
}

func (x Item) GetValue() (any, error) {
	return x.Token.GetValue()
}

func (x Item) GetTokenValue() (any, error) {
	return x.Token.GetValue()
}

func (x Item) GetRawPath() RawPath {
	clone := make([]PathItem, len(x.path))
	copy(clone, x.path)
	return clone
}

func (p RawPath) String() string {
	return fmt.Sprint(p)
}
func (p RawPath) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	for _, item := range p {
		fz.Printf("%s", item)
	}
}

func (p PathItem) IsArray() bool {
	return p.Token.typ == TokenArrayStart
}
func (p PathItem) IsObject() bool {
	return p.Token.typ == TokenObjectStart
}

// Value returns the value of the path item. If the item is inside an array, it returns the index. If the item is inside an object, it returns the key.
func (p PathItem) Value() any {
	switch p.Token.typ {
	case TokenArrayStart:
		return p.Index
	case TokenObjectStart:
		v, _ := p.Token.GetString()
		return v
	default:
		return nil
	}
}

func (p PathItem) String() string {
	return fmt.Sprint(p)
}
func (p PathItem) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	switch {
	case p.IsArray():
		fz.Printf("[%d]", p.Index)
	case p.IsObject():
		fz.Printf(".%s", p.Key)
	default:
		// root, empty path
	}
}
