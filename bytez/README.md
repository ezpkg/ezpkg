<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/bytez

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/bytez)](https://pkg.go.dev/ezpkg.io/bytez)
[![GitHub License](https://img.shields.io/github/license/ezpkg/bytez)](https://github.com/ezpkg/bytez/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/bytez?label=version)](https://github.com/ezpkg/bytez/tags)

Package bytez provides utilities for working with byte slices. It aims to extend the standard library `bytes` package with additional functionality.

## Installation

```sh
go get -u ezpkg.io/bytez@v0.0.4
```

## Examples

#### bytez.Buffer

The stdlib `bytes.Buffer` provides many functions that always return nil error. They have their counterparts in `bytez.Buffer` that eliminate the need of error handling.

```go
// stdlib: bytes.Buffer
_, err = b.WriteString()
if err != nil {
	return err
}
_, err = fmt.Fprintf(&b, "Hello, %s!", "world")
if err != nil {
	return err
}

// ezpkg.io/bytez.Buffer
b.WriteStringZ()
b.Printf("Hello, %s!", "world")
```

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
