package errorz_test

import (
	"fmt"
	"testing"

	"github.com/ezpkg/errorz"
)

type MockError struct{}

func (*MockError) Error() string { return "mock" }

func TestErrors(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var err error
		errorz.Append(&err, nil)
		errorz.Append(&err, (*MockError)(nil))
		errorz.Appendf(&err, nil, "empty error")

		assert(t, err == nil).Errorf("❌ expect err == nil")
	})
	t.Run("one", func(t *testing.T) {
		var err error
		errorz.Append(&err, (*MockError)(nil))
		errorz.Append(&err, fmt.Errorf("error/foo"))
		errorz.Append(&err, nil)

		errs := errorz.GetErrors(err)
		assert(t, len(errs) == 1).Errorf("❌ expect len(errs) == 1")
		assert(t, errs[0].Error() == "error/foo")
	})
	t.Run("many", func(t *testing.T) {
		var err1 error
		errorz.Append(&err1, fmt.Errorf("error/foo"))
		errorz.Append(&err1, errorz.Errorf("error/bar"))

		var err2 error
		errorz.Append(&err2, err1)
		errorz.Append(&err2, errorz.Errorf("error/baz"))

		errs1, errs2 := errorz.GetErrors(err1), errorz.GetErrors(err2)
		assert(t, len(errs1) == 2).Errorf("❌ expect len(errs1) == 2")
		assert(t, len(errs2) == 3).Errorf("❌ expect len(errs2) == 3")
	})
	t.Run("format", func(t *testing.T) {
		stackPlus := `
github.com/ezpkg/errorz_test.TestErrors.func█.█
	/Users/i/ws/ezpkg/ezpkg/errorz/multierr_test.go:██
testing.tRunner
	/usr/local/go/src/testing/testing.go:████
`
		stackAlt := `
github.com/ezpkg/errorz_test/multierr_test.go:██ · TestErrors.func█.█
testing/testing.go:████ · tRunner
`

		t.Run("one", func(t *testing.T) {
			var err error
			errorz.Append(&err, fmt.Errorf("error/foo"))

			t.Run("%v", func(t *testing.T) {
				str := fmt.Sprintf("%v", err)
				assertEqual(t, str, `(1 error) error/foo`)
			})
			t.Run("%+v", func(t *testing.T) {
				str := fmt.Sprintf("%+v", err)
				assertEqual(t, str, `
1 error occurred:
	* error/foo
`+stackPlus)
			})
			t.Run("%#v", func(t *testing.T) {
				str := fmt.Sprintf("%#v", err)
				assertEqual(t, str, `
1 error occurred:
	* error/foo
`+stackAlt)
			})
		})
		t.Run("two", func(t *testing.T) {
			var err error
			errorz.Append(&err, fmt.Errorf("error/foo"))
			errorz.Append(&err, fmt.Errorf("error/bar"))

			t.Run("%v", func(t *testing.T) {
				str := fmt.Sprintf("%v", err)
				assertEqual(t, str, `(2 errors) error/foo ; error/bar`)
			})
			t.Run("%+v", func(t *testing.T) {
				str := fmt.Sprintf("%+v", err)
				assertEqual(t, str, `
2 errors occurred:
	* error/foo
	* error/bar
`+stackPlus)
			})
			t.Run("%#v", func(t *testing.T) {
				str := fmt.Sprintf("%#v", err)
				assertEqual(t, str, `
2 errors occurred:
	* error/foo
	* error/bar
`+stackAlt)
			})
		})
	})
	t.Run("validate", func(t *testing.T) {
		var err error
		errorz.Validatef(&err, true, "ok, pass")
		errorz.Validatef(&err, false, "error/foo")
		errorz.ValidateX(&err, 42, false, "error/42")

		errs := errorz.GetErrors(err)
		assert(t, len(errs) == 2).Errorf("❌ expect len(errs) == 2")

		t.Run("format", func(t *testing.T) {
			str := fmt.Sprintf("%#v", err)
			assertEqual(t, str, `
2 errors occurred:
	* error/foo
	* error/42
github.com/ezpkg/errorz_test/multierr_test.go:███ · TestErrors.func█
testing/testing.go:████ · tRunner
`)
		})
	})
}
