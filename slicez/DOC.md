Package [slicez](https://pkg.go.dev/ezpkg.io/slicez) extends the standard library [slices](https://pkg.go.dev/slices) with additional functions.

## Examples

```go
codes := slicez.MapFunc([]int{1, 2, 3}, func(i int) string {
    return fmt.Sprintf("CODE(%d)", i)
})
fmt.Println(codes)
```
