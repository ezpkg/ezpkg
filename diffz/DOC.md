Package diffz provides functions for comparing and displaying differences between two strings. It's based on [kylelemons/godebug](https://github.com/kylelemons/godebug) and [sergi/go-diff](https://github.com/sergi/go-diff). It provides additional features of ignoring spaces and supporting placeholders.

## Examples

```go
// diff by char
left, right := "onetwo thr33four five", "onetwothree fourfive"
diffs := diffz.IgnoreSpace().DiffByChar(left, right)
fmt.Println(diffz.Format(diffs))

// diff by line
left, right := "one\ntwo\nthree\nfour", "one\ntwo\nfour"
diffs := diffz.IgnoreSpace().DiffByLine(left, right)
fmt.Println(diffz.Format(diffs))

// placeholder is useful for comparing tests with uuid or random values
diffs := diffz.Placeholder().AndIgnoreSpace().DiffByLine(left, right)
left := "id: ████\ncode: AF███\nname: Alice\n"
right := "id: 1234\ncode: AF123\nname: Baby\n"
fmt.Println(diffz.Format(diffs))
```

## Similar Packages

This package is based on these packages:

- [github.com/kylelemons/godebug](https://github.com/kylelemons/godebug)
- [github.com/sergi/go-diff](https://github.com/sergi/go-diff)
