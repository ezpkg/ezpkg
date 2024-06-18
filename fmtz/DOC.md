Package fmtz extends the standard library [fmt](https://pkg.go.dev/fmt) with additional functions.

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
