Package [iterz](https://pkg.go.dev/ezpkg.io/iterz) extends the standard library [iter](https://pkg.go.dev/iter) with additional functions.

## Features

Currently, it provides the following functions:

- `Nil[V]() iter.Seq[V]`: returns an iterator that yields nothing.
- `Nil2[K, V]() iter.Seq2[K, V]`: returns an iterator that yields nothing.

```go
// Nil return an iter.Seq that yields nothing.
func Nil[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {}
}

// Nil2 return an iter.Seq2 that yields nothing.
func Nil2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {}
}
```
