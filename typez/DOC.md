# ezpkg.io/typez

Package [typez](https://pkg.go.dev/ezpkg.io/typez) provides generic functions for working with types.

## Examples

```go
typez.In(1, 1, 2, 3)    // true
typez.In("A", "B", "C") // false

type A struct{X int}
typez.Coalesce(0, 1, 2, 3) // 1
typez.Coalesce(nil, &A{10}, &A{20}) // &A{10}
```
