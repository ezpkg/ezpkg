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
go get -u ezpkg.io/conveyz@v0.0.9
```

## Examples

```go
func Test(t *testing.T) {
	Ω := conveyz.GomegaExpect
	conveyz.Convey("Start", t, func() {
		s := "[0]"
		defer func() { fmt.Printf("\n%s\n", s) }()

		add := func(part string) {
			s = examples.AppendStr(s, part)
		}

		conveyz.Convey("Test 1:", func() {
			add(" → [1]")
			Ω(s).To(gomega.Equal("[0] → [1]"))

			conveyz.Convey("Test 1.1:", func() {
				add(" → [1.1]")
				Ω(s).To(gomega.Equal("[0] → [1] → [1.1]"))
			})
			conveyz.Convey("Test 1.2:", func() {
				add(" → [1.2]")
				Ω(s).To(gomega.Equal("[0] → [1] → [1.2]"))
			})
		})
		conveyz.Convey("Test 2:", func() {
			add(" → [2]")
			Ω(s).To(gomega.Equal("[0] → [2]"))

			conveyz.Convey("Test 2.1:", func() {
				add(" → [2.1]")
				Ω(s).To(gomega.Equal("[0] → [2] → [2.1]"))
			})
			conveyz.Convey("Test 2.2:", func() {
				add(" → [2.2]")
				Ω(s).To(gomega.Equal("[0] → [2] → [2.2]"))
			})
		})
		conveyz.SkipConvey("failure message", func() {
			// 👆 change SkipConvey to Convey to see failure messages

			conveyz.Convey("fail", func() {
				//  Expected
				//      <string>: [0] → [2]
				//  to equal
				//      <string>: this test will fail
				Ω(s).To(gomega.Equal("this test will fail"))
			})
			conveyz.Convey("UNEXPECTED ERROR", func() {
				// UNEXPECTED ERROR: Refusing to compare <nil> to <nil>.
				//  Be explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.
				Ω(nil).To(gomega.Equal(nil))
			})
		})
	})
}
```

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
