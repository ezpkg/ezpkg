<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/fmtz

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/fmtz)](https://pkg.go.dev/ezpkg.io/fmtz)
[![GitHub License](https://img.shields.io/github/license/ezpkg/fmtz)](https://github.com/ezpkg/fmtz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/fmtz?label=version)](https://github.com/ezpkg/fmtz/tags)

Package fmtz extends the standard library [fmt](https://pkg.go.dev/fmt) with additional functions.

## Installation

```sh
go get -u ezpkg.io/fmtz@v0.0.6
```

## Examples

#### fmt.State

The stdlib `fmt.State` provides many functions that always return nil error. They have their counterparts as `fmtz.State` that eliminate the need of error handling. There is also `fmtz.MustState` that panics on error, which is useful when other types implement `fmt.State` that may return non-nil error.

```go
type Code struct {
    Char   rune
    Number int
}

func (c Code) Format(s0 fmt.State, r rune) {
    s := fmtz.WrapState(s0)
    s.WriteRuneZ(c.Char)
    s.Print(c.Number)
}

func main() {
    c := Code{'Ω', 123}
    fmt.Printf("%v", c)
}
```

#### FormatMsgArgs

`fmtz.FormatMsgArgs` is a helper function that formats a message with arguments. It is useful for using in logging and error messages.

```go
func validate(err error, msgAndArgs ...any) error {
    if err == nil {
		return nil
    }
	msg := fmtz.FormatMsgArgs(msgAndArgs...)
	return typez.If(msg == "", err, fmt.Errorf("%v: %w", msg, err))
}

func main() {
    someError := errors.New("something went wrong")
    err := validate(someError, "failed to do something foo=%v bar=%v", "10", "20")
    fmt.Println(err)
	// Output: failed to do something foo=10 bar=20: something went wrong
}
```

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
