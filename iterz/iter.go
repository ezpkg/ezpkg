package iterz // import "ezpkg.io/iterz"

import (
	"iter"
)

// Nil return an iter.Seq that yields nothing.
func Nil[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {}
}

// Nil2 return an iter.Seq2 that yields nothing.
func Nil2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {}
}
