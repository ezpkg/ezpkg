package mapz

import (
	"slices"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

func FromSlice[S ~[]E, E comparable, R any](s S, fn func(E) R) map[E]R {
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
