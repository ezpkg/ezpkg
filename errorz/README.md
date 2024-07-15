<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/errorz

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/errorz)](https://pkg.go.dev/ezpkg.io/errorz)
[![GitHub License](https://img.shields.io/github/license/ezpkg/errorz)](https://github.com/ezpkg/errorz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/errorz?label=version)](https://pkg.go.dev/ezpkg.io/errorz?tab=versions)

Package [errorz](https://pkg.go.dev/ezpkg.io/errorz) provides functions for dealing with errors, with stacktrace, validation, and multi-errors.

## Installation

```sh
go get -u ezpkg.io/errorz@v0.1.1
```

## Examples

### Must

```go
data := errorz.Must(os.ReadFile(path))
fmt.Printf("%s", data)

errorz.MustZ(os.WriteFile(path, data, 0644))
```

### Stacktrace

```go
err := fmt.Errorf("foo")
zErr := errorz.Wrapf(err, "bar/%v", "one")

fmt.Printf("%v\n", zErr)
// bar/one: foo

fmt.Printf("%+v\n", zErr)
// bar/one: foo
// ezpkg.io/errorz_test.TestError.func1
// /Users/i/ws/ezpkg/ezpkg/errorz/errorz_test.go:12
// testing.tRunner
// /usr/local/go/src/testing/testing.go:7890

fmt.Printf("%#v\n", zErr)
// bar/one: foo
// ezpkg.io/errorz_test/errorz_test.go:12 · TestError.func1
// testing/testing.go:7890 · tRunner
```

### No stacktrace

```go
zErr := errorz.NoStack().New("no stack")
fmt.Printf("%+v", zErr)
```

### Multi-errors

```go
var err error
errorz.AppendTo(&err, fmt.Errorf("foo"))
errorz.AppendTo(&err, nil)
fmt.Printf("%+v", err)

var err2 error
err2 = errorz.Append(err2, err)
fmt.Printf("%+v", err2)
```

## Similar Packages

- [github.com/pkg/errors](https://pkg.go.dev/github.com/pkg/errors)
- [github.com/hashicorp/go-multierror](https://github.com/hashicorp/go-multierror)
- [github.com/uber-go/multierr](https://github.com/uber-go/multierr)
- [tailscale.com/util/multierr](https://pkg.go.dev/tailscale.com/util/multierr)
- [sigs.k8s.io/cli-utils/pkg/multierror](https://pkg.go.dev/sigs.k8s.io/cli-utils/pkg/multierror)

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
