package jsoniter

import (
	"errors"
	"fmt"
	"strconv"

	"ezpkg.io/fmtz"
)

var (
	ErrTokenInvalid = errors.New("invalid token")
	ErrTokenEmpty   = errors.New("invalid empty token")
	ErrTokenString  = errors.New("invalid string token")
	ErrTokenNumber  = errors.New("invalid number token")
	ErrNumberNotInt = errors.New("number is not an integer")
	ErrTokenBool    = errors.New("invalid boolean token")
	ErrTokenType    = errors.New("invalid token type")
)

type Item struct {
	Level int      // level of indentation
	Index int      // index in the parent array or object
	Key   RawToken // optional object "key"
	Token RawToken // [ or { or } or ] or , or value

	path RawPath
}

// Path is a slice of values. The values are the keys of objects (string) and the indexes of arrays (int).
type Path []any

type RawPath []PathItem

type PathItem struct {
	Index int      // array index or object index
	Key   RawToken // object key
	Token RawToken // [ or { or } or ]
}

func (x Item) String() string {
	return x.Token.String()
}

func (x Item) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	switch {
	case f.Flag('+'):
		fz.Printf("%v → %s", x.GetRawPath(), x.Token)

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

// GetPath returns the path of the item as a slice of values. The values are the keys of objects (string) and the indexes of arrays (int).
func (x Item) GetPath() Path {
	xPath := x.path[1:]
	path := make([]any, len(xPath))
	for i, item := range xPath {
		path[i] = item.Value()
	}
	return path
}

// GetRawPath returns the path of the item as a slice of PathItem.
// IMPORTANT: The result slice should not be modified.
func (x Item) GetRawPath() RawPath {
	return x.path[1:]
}

// GetPathString returns the path of the item as a string "0.key.1".
func (x Item) GetPathString() string {
	return fmt.Sprint(x.GetRawPath())
}

// GetAltPathString returns the path of the item as a string `[0].key[1]`.
func (x Item) GetAltPathString() string {
	return fmt.Sprintf("%+v", x.GetRawPath())
}

// Path returns the path of the item as a string. Default to 0.key.1 or "%+v" to format as [0]."key"[1]
func (p Path) String() string {
	return fmt.Sprint(p)
}

// Path returns the path of the item as a string. Default to 0.key.1 or "%+v" to format as [0]."key"[1]
func (p Path) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	plus := f.Flag('+')
	for i, item := range p {
		switch v := item.(type) {
		case int:
			if plus {
				fz.Printf("[%d]", v)
			} else {
				if i > 0 {
					fz.WriteByteZ('.')
				}
				fz.Printf("%d", v)
			}
		case string:
			if plus {
				fz.WriteByteZ('.')
			} else {
				if i > 0 {
					fz.WriteByteZ('.')
				}
			}
			if needQuote(v) {
				fz.Printf("%q", v)
			} else {
				fz.Printf("%s", v)
			}

		default:
			fz.Printf("!<%T>", v)
		}
	}
}

// NewRawPath returns a new RawPath. Use int for array index and string for object key. Keys should be quoted for better performance (optional).
func NewRawPath(parts ...any) (path RawPath, err error) {
	path = make(RawPath, len(parts))
	for i, part := range parts {
		path[i], err = NewPathItem(part)
		if err != nil {
			return nil, err
		}
	}
	return path, nil
}

// Match returns true if the path is equal to the given path, by raw values (int for array index, string for object key). Examples:
//
//	path.Match([]any{1, "key", 2})
func (p RawPath) Match(X ...any) bool {
	if len(X) == 1 {
		switch slice := X[0].(type) {
		case []any:
			return p.Match(slice...)
		case RawPath:
			return match(p, slice)
		default:
			// continue
		}
	}
	return match(p, X)
}

// Last returns the last item of the path.
func (p RawPath) Last() PathItem {
	if len(p) == 0 {
		return PathItem{}
	}
	return p[len(p)-1]
}

func match[T any](p RawPath, X []T) bool {
	if len(p) != len(X) {
		return false
	}
	for i, item := range X {
		if !p[i].Match(item) {
			return false
		}
	}
	return true
}

// ContainsRaw returns true if the path contains the given path, by raw values (int for array index, string for object key). Examples:
//
//	path.Contains([]any{1, "key", 2})
func (p RawPath) Contains(X ...any) bool {
	if len(p) < len(X) {
		return false
	}
	for i := 0; i <= len(p)-len(X); i++ {
		match := true
		for j := range X {
			if !p[i+j].Match(X[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// String returns the path of the item as a string. Default to 0.key.1 or "%+v" to format as [0]."key"[1]
func (p RawPath) String() string {
	return fmt.Sprint(p)
}

// Format formats the path as a string. Default to 0.key.1 or "%+v" to format as [0]."key"[1]
func (p RawPath) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	for i, item := range p {
		if f.Flag('+') {
			item.Format(f, c)
		} else {
			if i > 0 {
				fz.WriteByteZ('.')
			}
			item.Format(f, c)
		}
	}
}

// NewPathItem returns a new PathItem. Use int for array index and string for object key. Keys should be quoted for better performance (optional).
func NewPathItem(x any) (pi PathItem, err error) {
	switch x := intOrStr(x).(type) {
	case int:
		if x < 0 {
			return pi, fmt.Errorf("invalid path index: %d", x)
		}
		pi.Token = tokenArrayOpen
		pi.Index = x
		return pi, nil

	case string:
		pi.Token = tokenObjectOpen
		if len(x) == 0 {
			return pi, ErrTokenEmpty
		}
		if x[0] == '"' {
			pi.Key, err = NewRawToken([]byte(x))
		} else {
			pi.Key = RawToken{
				typ: TokenString,
				raw: quoteString(x),
			}
		}
		return pi, err

	case []byte:
		pi.Token = tokenObjectOpen
		if len(x) == 0 {
			return pi, ErrTokenEmpty
		}
		if x[0] == '"' {
			pi.Key, err = NewRawToken(x)
		} else {
			pi.Key = RawToken{
				typ: TokenString,
				raw: quoteString(string(x)),
			}
		}
		return pi, err

	default:
		return pi, fmt.Errorf("invalid path type: %T", x)
	}
}

// IsArray returns true if the path item is inside an array.
func (p PathItem) IsArray() bool {
	return p.Token.typ == TokenArrayOpen
}

// IsObject returns true if the path item is inside an object.
func (p PathItem) IsObject() bool {
	return p.Token.typ == TokenObjectOpen
}

// IsArrayIndex returns true if the path item is inside an array and the index matches the given index.
func (p PathItem) IsArrayIndex(idx int) bool {
	return p.Token.typ == TokenArrayOpen && p.Index == idx
}

// IsObjectRawKey returns true if the path item is inside an object and the key matches the given unquote key.
func (p PathItem) IsObjectRawKey(rawKey string) bool {
	return p.Token.typ == TokenObjectOpen && string(p.Key.raw) == rawKey
}

// IsObjectKey returns true if the path item is inside an object and the key matches the given key.
func (p PathItem) IsObjectKey(key string) bool {
	if p.Token.typ == TokenObjectOpen {
		pKey, err := p.Key.GetString()
		return err == nil && pKey == key
	}
	return false
}

// ArrayIndex returns the array index of the path item. It returns the index as an int if the item is inside an array.
func (p PathItem) ArrayIndex() (int, bool) {
	if p.Token.typ == TokenArrayOpen {
		return p.Index, true
	}
	return 0, false
}

// ObjectKey returns the object key of the path item. It returns the key as a string if the item is inside an object.
func (p PathItem) ObjectKey() (string, bool) {
	if p.Token.typ == TokenObjectOpen {
		if str, err := p.Key.GetString(); err == nil {
			return str, true
		}
	}
	return "", false
}

// Match returns true if the path item matches the given value. The value must be an int for array index or a string|[]byte for object key.
func (p PathItem) Match(x any) bool {
	switch x := intOrStr(x).(type) {
	case int:
		return p.IsArrayIndex(x)
	case string:
		return p.IsObjectKey(x)
	case []byte:
		return p.IsObjectKey(string(x))
	case PathItem:
		if p.Token.typ != x.Token.typ {
			return false
		}
		switch p.Token.typ {
		case TokenArrayOpen:
			return p.Index == x.Index
		case TokenObjectOpen:
			return p.Key.Equal(x.Key)
		default:
			return false
		}
	default:
		return false
	}
}

// Value returns the value of the path item. If the item is inside an array, it returns the index. If the item is inside an object, it returns the key.
func (p PathItem) Value() any {
	switch p.Token.typ {
	case TokenArrayOpen:
		return p.Index
	case TokenObjectOpen:
		v, _ := p.Key.GetString()
		return v
	default:
		return nil
	}
}

// String returns the string representation of the path item. "0" for array, "key" for object.
func (p PathItem) String() string {
	switch {
	case p.IsArray():
		return strconv.Itoa(p.Index)
	case p.IsObject():
		rawKey := p.Key.Raw()
		if canSimplyUnquote(rawKey) {
			rawKey = rawKey[1 : len(rawKey)-1]
		}
		return string(rawKey)
	default:
		return ""
	}
}

// Format formats the path item as a string.
// Use "%+v" to format as "[0]" for array, ".key" for object.
func (p PathItem) Format(f fmt.State, c rune) {
	fz := fmtz.WrapState(f)
	switch {
	case p.IsArray():
		if f.Flag('+') {
			fz.Printf("[%d]", p.Index)
		} else {
			fz.Printf("%d", p.Index)
		}
	case p.IsObject():
		rawKey := p.Key.Raw()
		if canSimplyUnquote(rawKey) {
			rawKey = rawKey[1 : len(rawKey)-1]
		}
		if f.Flag('+') {
			fz.Printf(".%s", rawKey)
		} else {
			fz.Printf("%s", rawKey)
		}
	default:
		// root, empty path
	}
}
