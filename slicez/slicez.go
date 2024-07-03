package slicez // import "ezpkg.io/slicez"

import (
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
