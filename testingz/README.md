<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/testingz

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/testingz)](https://pkg.go.dev/ezpkg.io/testingz)
[![GitHub License](https://img.shields.io/github/license/ezpkg/testingz)](https://github.com/ezpkg/testingz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/testingz?label=version)](https://github.com/ezpkg/testingz/tags)

Packages testingz provides utilities for testing. Support ignoring spaces and using placeholders.

## Installation

```sh
go get -u ezpkg.io/testingz@v0.0.4
```

## Examples

```go
formatted, isDiff := testingz.DiffByChar("code: 123A", "code: ███A")
// isDiff: true

formatted, isDiff := testingz.DiffByCharZ("code: 123A", "code: ███A")
// isDiff: false

// placeholder is useful for comparing tests with uuid or random values
formatted, isDiff := testingz.DiffByLineZ(left, right)
left := "id: ████\ncode: AF███\nname: Alice\n"
right := "id: 1234\ncode: AF123\nname: Alice\n"
// isDiff: false
```

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
