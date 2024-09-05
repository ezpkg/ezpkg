package typez // import "ezpkg.io/typez"

import (
	"reflect"
)

// --- pointer --- //

func Ptr[T any](v T) *T {
	return &v
}

func Deptr[T any](v *T) (out T) {
	if v != nil {
		return *v
	}
	return out
}

// --- branch --- //

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

// --- comparable --- //

func In[T comparable](item T, list ...T) bool {
	for _, x := range list {
		if item == x {
			return true
		}
	}
	return false
}

func Coalesce[T comparable](list ...T) T {
	var zero T
	for _, item := range list {
		if item != zero {
			return item
		}
	}
	return zero
}

func CoalesceX[T any](list ...T) (out T) {
	for _, elem := range list {
		if !IsNil(elem) {
			return elem
		}
	}
	return out
}

// IsNil reports whether v is nil. Unlike reflect.IsNil(), it won't panic.
//
// The Go team decided not to add "zero" to the language.
// https://github.com/golang/go/issues/61372
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

// IsZero reports whether v is zero. Unlike reflect.IsZero(), it won't panic.
//
// The Go team decided not to add "zero" to the language.
// https://github.com/golang/go/issues/61372
func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}
