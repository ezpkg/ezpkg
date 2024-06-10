Package stacktracez provides functions for getting stack trace for using in errors and logs.

## Examples

```go
stack := stacktracez.StackTrace()

fmt.Printf("%+v", stack)
// ezpkg.io/stacktracez.TestStackTrace.func1
// /Users/i/ws/ezpkg/ezpkg/stacktracez/stacktracez_test.go:12
// ezpkg.io/stacktracez.TestStackTrace
// /Users/i/ws/ezpkg/ezpkg/stacktracez/stacktracez_test.go:15
// testing.tRunner
// /usr/local/go/src/testing/testing.go:1689

fmt.Printf("%v", stack)
// ezpkg.io/stacktracez/stacktracez_test.go:12 · TestStackTrace.func1
// ezpkg.io/stacktracez/stacktracez_test.go:15 · TestStackTrace
// testing/testing.go:1689 · tRunner
```
