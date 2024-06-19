Package [unsafez](https://pkg.go.dev/ezpkg.io/unsafez) provides unsafe functions for working with unsafe.Pointer.

## Examples

```go
data := errorz.Must(os.ReadFile()))
str := unsafez.BytesToString(data)
fmt.Println(str)
```
