Package colorz provides utilities for working with colors in terminal.

## Examples

### colorz.Color

```go
colorz.Yellow.Wrap("this is yellow")
colorz.Red.Printf("error: %s", "something went wrong")
fmt.Printf("roses are %sred%s and violets are %sblue%s\n", colorz.Red,colorz.Reset, colorz.Green, colorz.Reset)
```

## Similar Packages

- [github.com/fatih/color](https://github.com/fatih/color)
- [github.com/gookit/color](https://github.com/gookit/color)
