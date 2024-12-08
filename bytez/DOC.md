# ezpkg.io/bytez

Package [bytez](https://pkg.go.dev/ezpkg.io/bytez) provides utilities for working with byte slices. It aims to extend the standard library [bytes](https://pkg.go.dev/bytes) package with additional functionality.

## Examples

### bytez.Buffer

The stdlib `bytes.Buffer` provides many functions that always return nil error. They have their counterparts in `bytez.Buffer` that eliminate the need of error handling.

```go
// stdlib: bytes.Buffer
_, err = b.WriteString()
if err != nil {
	return err
}
_, err = fmt.Fprintf(&b, "Hello, %s!", "world")
if err != nil {
	return err
}

// ezpkg.io/bytez.Buffer
b.WriteStringZ()
b.Printf("Hello, %s!", "world")
```
