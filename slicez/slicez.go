package slicez // import "ezpkg.io/slicez"

import (
	"cmp"
	"slices"

	"ezpkg.io/typez"
)

// GetX returns the element at index n. If n is negative, it returns from the end.
func GetX[S ~[]E, E any](s S, n int) (out E, ok bool) {
	if n < 0 {
		n = len(s) + n
	}
	if n >= 0 && n < len(s) {
		return s[n], true
	}
	return out, false
}

// Get returns the element at index n. If n is negative, it returns from the end.
func Get[S ~[]E, E any](s S, n int) (out E) {
	out, _ = GetX(s, n)
	return out
}

// GetOr returns the element at index n. If n is negative, it returns from the end. If n is out of range, it returns the fallback value.
func GetOr[S ~[]E, E any](s S, n int, fallback E) (out E) {
	out, ok := GetX(s, n)
	return typez.If(ok, out, fallback)
}

// GetOrFunc returns the element at index n. If n is negative, it returns from the end. If n is out of range, it returns the result of the fallback function.
func GetOrFunc[S ~[]E, E any](s S, n int, fallback func() E) (out E) {
	out, ok := GetX(s, n)
	if ok {
		return out
	} else {
		return fallback()
	}
}

// GetFunc returns the element that satisfies the function. To get the index, use slices.Index or slices.IndexFunc.
func GetFunc[S ~[]E, E any](s S, fn func(E) bool) (out E) {
	for _, item := range s {
		if fn(item) {
			return item
		}
	}
	return out
}

// Exists returns true if the element exists in the slice. To get the index, use slices.Index or slices.IndexFunc.
func Exists[S ~[]E, E comparable](s S, item E) bool {
	return ExistsFunc(s, func(x E) bool {
		return x == item
	})
}

// Exists returns true if the element exists in the slice. To get the index, use slices.Index or slices.IndexFunc.
func ExistsFunc[S ~[]E, E any](s S, fn func(E) bool) bool {
	for _, item := range s {
		if fn(item) {
			return true
		}
	}
	return false
}

// First returns the first element of the slice.
func First[S ~[]E, E any](s S) (out E) {
	if len(s) > 0 {
		return s[0]
	}
	return out
}

// FirstN returns the first n elements of the slice.
func FirstN[S ~[]E, E any](s S, n int) (out []E) {
	if n < 0 {
		return LastN(s, -n)
	}
	n = min(n, len(s))
	return s[:n]
}

// FirstFunc returns the first element of the slice that satisfies the function.
func FirstFunc[S ~[]E, E any](s S, fn func(E) bool) (out E) {
	for _, item := range s {
		if fn(item) {
			return item
		}
	}
	return out
}

// Last returns the last element of the slice.
func Last[S ~[]E, E any](s S) (out E) {
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return out
}

// LastN returns the last n elements of the slice.
func LastN[S ~[]E, E any](s S, n int) (out []E) {
	if n < 0 {
		return FirstN(s, -n)
	}
	n = min(n, len(s))
	return s[len(s)-n:]
}

// LastFunc returns the last element of the slice that satisfies the function.
func LastFunc[S ~[]E, E any](s S, fn func(E) bool) (out E) {
	for i := len(s) - 1; i >= 0; i-- {
		if fn(s[i]) {
			return s[i]
		}
	}
	return out
}

// MapFunc returns a new slice with elements mapped to a new type.
func MapFunc[S ~[]E, E, R any](s S, fn func(E) R) []R {
	result := make([]R, len(s))
	for i, item := range s {
		result[i] = fn(item)
	}
	return result
}

// FilterFunc returns a new slice with elements that satisfy the function.
func FilterFunc[S ~[]E, E any](s S, fn func(E) bool) (outs []E) {
	for _, item := range s {
		if fn(item) {
			outs = append(outs, item)
		}
	}
	return outs
}

// MapFilterFunc (or FilterMapFunc) returns a new slice with elements that satisfy the function, and maps them to a new type.
func MapFilterFunc[S ~[]E, E, R any](s S, fn func(E) (R, bool)) (outs []R) {
	for _, item := range s {
		out, ok := fn(item)
		if ok {
			outs = append(outs, out)
		}
	}
	return outs
}

// FilterMapFunc (or FilterMapFunc) returns a new slice with elements that satisfy the function, and maps them to a new type.
func FilterMapFunc[S ~[]E, E any](s S, fn func(E) (E, bool)) (outs []E) {
	return MapFilterFunc(s, fn)
}

// Concat combine multiple slices into a new slice.
func Concat[S ~[]E, E any](slices ...S) []E {
	return AppendSlice(nil, slices...)
}

// AppendSlice appends multiple slices to the first slice.
func AppendSlice[S ~[]E, E any](s S, slices ...S) []E {
	N := len(s)
	for _, slice := range slices {
		N += len(slice)
	}
	outs := make([]E, 0, N)
	outs = append(outs, s...)
	for _, slice := range slices {
		outs = append(outs, slice...)
	}
	return outs
}

// AppendTo appends multiple elements to a slice pointer.
func AppendTo[S ~*[]E, E any](s S, items ...E) []E {
	*s = append(*s, items...)
	return *s
}

// AppendSliceTo appends multiple slices to a slice pointer.
func AppendSliceTo[S ~*[]E, E any](s S, slices ...[]E) []E {
	*s = AppendSlice(*s, slices...)
	return *s
}

// Prepend appends multiple elements to the beginning of a slice.
func Prepend[S ~[]E, E any](s S, items ...E) []E {
	return append(items, s...)
}

// PrependTo appends multiple elements to the beginning of a slice pointer.
func PrependTo[S ~*[]E, E any](s S, items ...E) []E {
	*s = append(items, *s...)
	return *s
}

// IsUnique returns true if all elements in the slice are unique.
func IsUnique[S ~[]E, E comparable](s S) bool {
	m := make(map[E]struct{})
	for _, item := range s {
		if _, ok := m[item]; ok {
			return false
		}
		m[item] = struct{}{}
	}
	return true
}

// IsUniqueFunc returns true if all elements in the slice are unique based on the function.
func IsUniqueFunc[S ~[]E, E any, K comparable](s S, fn func(E) K) bool {
	m := make(map[K]struct{})
	for _, item := range s {
		key := fn(item)
		if _, ok := m[key]; ok {
			return false
		}
		m[key] = struct{}{}
	}
	return true
}

// Unique returns a slice with unique elements. If the slice is already unique, it returns the original slice.
func Unique[S ~[]E, E comparable](s S) (out []E) {
	m := make(map[E]struct{})
	for i, item := range s {
		if _, exist := m[item]; !exist {
			m[item] = struct{}{}
		} else {
			out = s[:i:i]
			break // use slow fast
		}
	}
	if out == nil { // all unique
		return s
	}
	// slow path, check remaining
	for _, item := range s[len(out)+1:] {
		if _, exist := m[item]; !exist {
			m[item] = struct{}{}
			out = append(out, item)
		}
	}
	return out
}

// UniqueFunc returns a new slice with unique elements.
func UniqueFunc[S ~[]E, E any, K comparable](s S, fn func(E) K) (out []E) {
	m := make(map[K]struct{})
	for i, item := range s {
		key := fn(item)
		if _, exist := m[key]; !exist {
			m[key] = struct{}{}
		} else {
			out = s[:i:i]
			break // use slow fast
		}
	}
	if out == nil { // all unique
		return s
	}
	// slow path, check remaining
	for _, item := range s[len(out):] {
		key := fn(item)
		if _, exist := m[key]; !exist {
			m[key] = struct{}{}
			out = append(out, item)
		}
	}
	return out
}

// SortUnique modified content of the slice, and return the modified slice with unique elements, sorted.
func SortUnique[S ~[]E, E cmp.Ordered](s S) (out []E) {
	slices.Sort(s)
	return slices.Compact(s)
}

// SortUniqueFunc modified content of the slice, and return the modified slice with unique elements, sorted based on the function.
func SortUniqueFunc[S ~[]E, E any, K cmp.Ordered](s S, fn func(E) K) (out []E) {
	slices.SortFunc(s, func(a, b E) int {
		return cmp.Compare(fn(a), fn(b))
	})
	return slices.CompactFunc(s, func(a, b E) bool {
		return fn(a) == fn(b)
	})
}
