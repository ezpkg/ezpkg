<p align="center">
<a href="https://ezpkg.io">
<img alt="gopherz" src="https://ezpkg.io/_/gopherz.png" style="width:420px">
</a>
</p>

# ezpkg.io/stringz

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ezpkg/stringz)](https://pkg.go.dev/github.com/ezpkg/stringz/v2)
[![GitHub License](https://img.shields.io/github/license/ezpkg/stringz)](https://github.com/ezpkg/stringz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/stringz?label=version)](https://github.com/ezpkg/stringz/tags)

Package stringz extends the standard library [strings](https://pkg.go.dev/strings) with additional functions.

## Installation

```sh
go get -u ezpkg.io/stringz@v0.0.2
```

## Examples

#### stringz.Builder

The stdlib `strings.Builder` provides many functions that always return nil error. They have their counterparts in `stringz.Builder` that eliminate the need of error handling.

```go
// stdlib: strings.Builder
_, err = b.WriteString()
if err != nil {
    return err
}
_, err = fmt.Fprintf(&b, "Hello, %s!", "world")
if err != nil {
    return err
}

// ezpkg.io/stringz.Builder
b.WriteStringZ()
b.Printf("Hello, %s!", "world")
```

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

<a href="https://olivernguyen.io"><img alt="olivernguyen.io" src="https://olivernguyen.io/_/badge.png" height="28px"></a>&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
