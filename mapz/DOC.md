Package mapz extends the package [golang.org/x/exp/maps](https://pkg.go.dev/golang.org/x/exp/maps) with additional functions.

## Examples

```go
mapCodes := mapz.FromSliceFunc([]int{1, 2, 3}, func(i int) string {
    return fmt.Sprintf("CODE(%d)", i)
})
fmt.Println(mapCodes)
```

## Similar Packages

This package is based on:

- [golang.org/x/exp/maps](https://pkg.go.dev/golang.org/x/exp/maps)
