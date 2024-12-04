<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/stringz

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/stringz)](https://pkg.go.dev/ezpkg.io/stringz)
[![GitHub License](https://img.shields.io/github/license/ezpkg/stringz)](https://github.com/ezpkg/stringz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/stringz?label=version)](https://pkg.go.dev/ezpkg.io/stringz?tab=versions)

Package [stringz](https://pkg.go.dev/ezpkg.io/stringz) extends the standard library [strings](https://pkg.go.dev/strings) with additional functions.

## Installation

```sh
go get -u ezpkg.io/stringz@v0.2.0
```

## Examples

### stringz.Builder

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

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
