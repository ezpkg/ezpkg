# ezpkg.io/fmtz

Package [fmtz](https://pkg.go.dev/ezpkg.io/fmtz) extends the standard library [fmt](https://pkg.go.dev/fmt) with additional functions.

## Examples

### fmt.State

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

### FormatMsgArgs

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
