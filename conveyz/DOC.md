# ezpkg.io/conveyz

Package [conveyz](https://pkg.go.dev/ezpkg.io/conveyz) extends the package [convey](https://pkg.go.dev/github.com/smartystreets/goconvey/convey) with additional functionality and make it work with [gomega](https://pkg.go.dev/github.com/onsi/gomega). See [the original blog post](https://olivernguyen.io/w/goconvey.gomega/).

## Features

- [Convey](https://pkg.go.dev/ezpkg.io/conveyz#Convey) functions to group tests. [SConvey](https://pkg.go.dev/ezpkg.io/conveyz#SConvey) to skip and [FConvey](https://pkg.go.dev/ezpkg.io/conveyz#FConvey) to focus on a test.
- Tests can be nested, and will be executed in nested order. In the example below, it will print:
```
[0] → [1] → [1.1]
[0] → [1] → [1.2]
[0] → [2] → [2.1]
[0] → [2] → [2.2]
```
  
### Differences to the original package

- This package [ezpkg.io/conveyz](https://pkg.go.dev/ezpkg.io/conveyz) uses [gomega](https://pkg.go.dev/github.com/onsi/gomega) for assertions. To me, gomega is more powerful and easier to use.
- With this package, we only need to set [FConvey](https://pkg.go.dev/ezpkg.io/conveyz#FConvey) at a single block level. The original package requires to set [FocusConvey](https://pkg.go.dev/github.com/smartystreets/goconvey/convey#FocusConvey) at every nested level.
- [FConvey](https://pkg.go.dev/ezpkg.io/conveyz#FConvey) and [SConvey](https://pkg.go.dev/ezpkg.io/conveyz#SConvey) will make `go test` fail. This is to avoid accidentally skipping tests. To skip tests but not make `go test` fail, use [SkipConveyAsTODO](https://pkg.go.dev/ezpkg.io/conveyz#SkipConveyAsTODO).
- Output looks better with more colors. Stacktrace is nicer. 

## Examples

See [conveyz/examples/conveyz_test.go](https://github.com/ezpkg/ezpkg/blob/main/conveyz/examples/conveyz_test.go) or [stringz/stringz_test.go](https://github.com/ezpkg/ezpkg/blob/main/stringz/stringz_test.go) for more examples.

```go
import (
	"fmt"
	"testing"

	. "ezpkg.io/conveyz"
)

func Test(t *testing.T) {
	Convey("Start", t, func() {
		s := "[0]"
		defer func() { fmt.Printf("\n%s\n", s) }()

		add := func(part string) {
			s = AppendStr(s, part)
		}

		Convey("Test 1:", func() {
			add(" → [1]")
			Ω(s).To(Equal("[0] → [1]"))

			Convey("Test 1.1:", func() {
				add(" → [1.1]")
				Ω(s).To(Equal("[0] → [1] → [1.1]"))
			})
			Convey("Test 1.2:", func() {
				add(" → [1.2]")
				Ω(s).To(Equal("[0] → [1] → [1.2]"))
			})
		})
        // 👇change to FConvey to focus on this block and all children
        // 👇change to SConvey to skip the block
		// 👇change to SkipConveyAsTODO to mark as TODO
		Convey("Test 2:", func() {
			add(" → [2]")
			Ω(s).To(Equal("[0] → [2]"))
			
			Convey("Test 2.1:", func() {
				add(" → [2.1]")
				Ω(s).To(Equal("[0] → [2] → [2.1]"))
			})
			Convey("Test 2.2:", func() {
				add(" → [2.2]")
				Ω(s).To(Equal("[0] → [2] → [2.2]"))
			})
		})
		SkipConvey("failure message", func() {
			// 👆change SkipConvey to Convey to see failure messages
			// 👆change SkipConvey to SkipConveyAsTODO to mark as TODO

			Convey("fail", func() {
				//  Expected
				//      <string>: [0] → [2]
				//  to equal
				//      <string>: this test will fail
				Ω(s).To(Equal("this test will fail"))
			})
			Convey("UNEXPECTED", func() {
				// UNEXPECTED: Refusing to compare <nil> to <nil>.
				//  Be explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.
				Ω(nil).To(Equal(nil))
			})
		})
	})
}

func AppendStr(s, part string) string {	return s + part }
func WillPanic() { panic("let's panic! 💥") }
func CallFunc(fn func()) { fn() }
```

### The output will look like this:

![](https://olivernguyen.io/w/ezpkg/_/cv1.png)
