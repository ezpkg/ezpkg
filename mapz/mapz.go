package mapz // import "ezpkg.io/mapz"

import (
	"slices"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"

	"ezpkg.io/typez"
)

func FromSliceFunc[S ~[]E, E comparable, R any](s S, fn func(E) R) map[E]R {
	result := make(map[E]R, len(s))
	for _, item := range s {
		result[item] = fn(item)
	}
	return result
}

func SortedKeys[M ~map[K]V, K constraints.Ordered, V any](m M) []K {
	keys := maps.Keys(m)
	slices.Sort(keys)
	return keys
}

func SortedKeysAndValues[M ~map[K]V, K constraints.Ordered, V any](m M) ([]K, []V) {
	keys := maps.Keys(m)
	slices.Sort(keys)
	values := make([]V, len(keys))
	for i, key := range keys {
		values[i] = m[key]
	}
	return keys, values
}

func Append[M ~map[K]V, K comparable, V any](m M, items map[K]V) map[K]V {
	if m == nil {
		m = make(map[K]V, len(items))
	}
	for k, v := range items {
		m[k] = v
	}
	return m
}

func Merge[M ~map[K]V, K comparable, V any](maps ...M) map[K]V {
	N := 0
	for _, x := range maps {
		N += len(x)
	}
	m := make(map[K]V, N)
	for _, x := range maps {
		for k, v := range x {
			m[k] = v
		}
	}
	return m
}

func Exists[M ~map[K]V, K comparable, V any](m M, key K) bool {
	_, ok := m[key]
	return ok
}

func ExistsFunc[M ~map[K]V, K comparable, V any](m M, fn func(K, V) bool) bool {
	for k, v := range m {
		if fn(k, v) {
			return true
		}
	}
	return false
}

func ExistsAll[M ~map[K]V, K comparable, V any](m M, keys ...K) bool {
	for _, key := range keys {
		if !Exists(m, key) {
			return false
		}
	}
	return true
}

func ExistsAny[M ~map[K]V, K comparable, V any](m M, keys ...K) bool {
	for _, key := range keys {
		if Exists(m, key) {
			return true
		}
	}
	return false
}

func FilterFunc[M ~map[K]V, K comparable, V any](m M, fn func(K, V) bool) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		if fn(k, v) {
			result[k] = v
		}
	}
	return result
}

func FilterKeys[M ~map[K]V, K comparable, V any](m M, keys ...K) map[K]V {
	result := make(map[K]V, len(m))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}
	return result
}

func FilterValues[M ~map[K]V, K, V comparable](m M, values ...V) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		if typez.In(v, values...) {
			result[k] = v
		}
	}
	return result
}
