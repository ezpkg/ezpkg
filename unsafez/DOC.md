Package unsafez provides unsafe functions for working with unsafe.Pointer.

## Examples

```go
data, err := os.ReadFile())
errorz.MustZ(err)
str := unsafez.BytesToString(data)
fmt.Println(str)
```
