<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/conveyz

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/conveyz)](https://pkg.go.dev/ezpkg.io/conveyz)
[![GitHub License](https://img.shields.io/github/license/ezpkg/conveyz)](https://github.com/ezpkg/conveyz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/conveyz?label=version)](https://pkg.go.dev/ezpkg.io/conveyz?tab=versions)

Package [conveyz](https://pkg.go.dev/ezpkg.io/conveyz) extends the package [convey](https://pkg.go.dev/github.com/smartystreets/goconvey/convey) with additional functionality and make it work with [gomega](https://pkg.go.dev/github.com/onsi/gomega). See [the original blog post](https://olivernguyen.io/w/goconvey.gomega/).

## Installation

```sh
go get -u ezpkg.io/conveyz@v0.1.0
```

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

	. "github.com/onsi/gomega"
	. "ezpkg.io/conveyz"
)

func Test(t *testing.T) {
	Ω := GomegaExpect // 👈 make Ω as alias for GomegaExpect
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

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
