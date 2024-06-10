Package stringz extends the standard library [strings](https://pkg.go.dev/strings) with additional functions.

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
