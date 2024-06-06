package typez

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
	for _, item := range list {
		if !reflect.ValueOf(item).IsNil() {
			return item
		}
	}
	return out
}
