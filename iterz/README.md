<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/iterz

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/iterz)](https://pkg.go.dev/ezpkg.io/iterz)
[![GitHub License](https://img.shields.io/github/license/ezpkg/iterz)](https://github.com/ezpkg/iterz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/iterz?label=version)](https://pkg.go.dev/ezpkg.io/iterz?tab=versions)

Package [iterz](https://pkg.go.dev/ezpkg.io/iterz) extends the standard library [iter](https://pkg.go.dev/iter) with additional functions.

## Installation

```sh
go get -u ezpkg.io/iterz@v0.2.0
```

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

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
