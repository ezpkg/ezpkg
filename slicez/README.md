<p align="center">
<a href="https://ezpkg.io">
<img alt="gopherz" src="https://ezpkg.io/_/gopherz.png" style="width:420px">
</a>
</p>

# ezpkg.io/slicez

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ezpkg/slicez)](https://pkg.go.dev/github.com/ezpkg/slicez/v2)
[![GitHub License](https://img.shields.io/github/license/ezpkg/slicez)](https://github.com/ezpkg/slicez/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/slicez?label=version)](https://github.com/ezpkg/slicez/tags)

Package slicez extends the standard library [slices](https://pkg.go.dev/slices) with additional functions.

## Installation

```sh
go get -u ezpkg.io/slicez@v0.0.2
```

## Examples

```go
codes := slicez.MapFunc([]int{1, 2, 3}, func(i int) string {
    return fmt.Sprintf("CODE(%d)", i)
})
fmt.Println(codes)
```

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

<a href="https://olivernguyen.io"><img alt="olivernguyen.io" src="https://olivernguyen.io/_/badge.png" height="28px"></a>&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
