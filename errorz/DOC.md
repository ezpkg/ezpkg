Package errorz provides functions for dealing with errors, with stacktrace, validation, and multi-errors.

## Examples

#### Must

```go
data := errorz.Must(os.ReadFile(path))
fmt.Printf("%s", data)

errorz.MustZ(os.WriteFile(path, data, 0644))
```

#### Stacktrace

```go
err := fmt.Errorf("foo")
zErr := errorz.Wrapf(err, "bar/%v", "one")

fmt.Printf("%v\n", zErr)
// bar/one: foo

fmt.Printf("%+v\n", zErr)
// bar/one: foo
// ezpkg.io/errorz_test.TestError.func1
// /Users/i/ws/ezpkg/ezpkg/errorz/errorz_test.go:12
// testing.tRunner
// /usr/local/go/src/testing/testing.go:7890

fmt.Printf("%#v\n", zErr)
// bar/one: foo
// ezpkg.io/errorz_test/errorz_test.go:12 · TestError.func1
// testing/testing.go:7890 · tRunner
```

**No stacktrace:**

```go
zErr := errorz.NoStack().New("no stack")
fmt.Printf("%+v", zErr)
```

#### Multi-errors

```go
var err error
errorz.AppendTo(&err, fmt.Errorf("foo"))
errorz.AppendTo(&err, nil)
fmt.Printf("%+v", err)

var err2 error
err2 = errorz.Append(err2, err)
fmt.Printf("%+v", err2)
```
