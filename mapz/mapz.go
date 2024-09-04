package mapz // import "ezpkg.io/mapz"

import (
	"slices"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
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

func Append[K comparable, V any](m map[K]V, items map[K]V) map[K]V {
	if m == nil {
		m = make(map[K]V, len(items))
	}
	for k, v := range items {
		m[k] = v
	}
	return m
}

func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
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
