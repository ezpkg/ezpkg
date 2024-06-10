<p align="center">
<a href="https://ezpkg.io">
<img alt="gopherz" src="https://ezpkg.io/_/gopherz.png" style="width:420px">
</a>
</p>

# ezpkg.io/colorz

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ezpkg/colorz)](https://pkg.go.dev/github.com/ezpkg/colorz/v2)
[![GitHub License](https://img.shields.io/github/license/ezpkg/colorz)](https://github.com/ezpkg/colorz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/colorz?label=version)](https://github.com/ezpkg/colorz/tags)

Package colorz provides utilities for working with colors in terminal.

## Installation

```sh
go get -u ezpkg.io/colorz@v0.0.1
```

## Examples

### colorz.Color

```go
colorz.Yellow.Wrap("this is yellow")
colorz.Red.Printf("error: %s", "something went wrong")
fmt.Printf("roses are %sred%s and violets are %sblue%s\n", colorz.Red,colorz.Reset, colorz.Green, colorz.Reset)
```

## Similar Packages

- [github.com/fatih/color](https://github.com/fatih/color)
- [github.com/gookit/color](https://github.com/gookit/color)

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

<a href="https://olivernguyen.io"><img alt="olivernguyen.io" src="https://olivernguyen.io/_/badge.png" height="28px"></a>&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
