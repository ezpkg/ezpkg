package errorz_test

import (
	"errors"
	"fmt"
	"testing"

	"ezpkg.io/errorz"
)

type MockError struct{}

func (*MockError) Error() string { return "mock" }

func TestErrors(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var err error
		errorz.AppendTo(&err, nil)
		errorz.AppendTo(&err, (*MockError)(nil))
		errorz.AppendTof(&err, nil, "empty error")

		assert(t, err == nil).Errorf("❌ expect err == nil")
	})
	t.Run("one", func(t *testing.T) {
		var err error
		errorz.AppendTo(&err, (*MockError)(nil))
		errorz.AppendTo(&err, fmt.Errorf("error/foo"))
		errorz.AppendTo(&err, nil)

		errs := errorz.GetErrors(err)
		assert(t, len(errs) == 1).Errorf("❌ expect len(errs) == 1")
		assert(t, errs[0].Error() == "error/foo")
	})
	t.Run("many", func(t *testing.T) {
		var err1 error
		errorz.AppendTo(&err1, fmt.Errorf("error/foo"))
		errorz.AppendTo(&err1, errorz.Errorf("error/bar"))

		var err2 error
		errorz.AppendTo(&err2, err1)
		errorz.AppendTo(&err2, errorz.Errorf("error/baz"))

		errs1, errs2 := errorz.GetErrors(err1), errorz.GetErrors(err2)
		assert(t, len(errs1) == 2).Errorf("❌ expect len(errs1) == 2")
		assert(t, len(errs2) == 3).Errorf("❌ expect len(errs2) == 3")
	})
	t.Run("format", func(t *testing.T) {
		stackPlus := `
ezpkg.io/errorz_test.TestErrors.func█.█
	/Users/i/ws/ezpkg/ezpkg/errorz/multierr_test.go:██
testing.tRunner
	/usr/local/go/src/testing/testing.go:████
`
		stackAlt := `
ezpkg.io/errorz_test/multierr_test.go:██ · TestErrors.func█.█
testing/testing.go:████ · tRunner
`

		t.Run("one", func(t *testing.T) {
			var err error
			errorz.AppendTo(&err, fmt.Errorf("error/foo"))

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
			errorz.AppendTo(&err, fmt.Errorf("error/foo"))
			errorz.AppendTo(&err, fmt.Errorf("error/bar"))

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
		errorz.ValidateTof(&err, true, "ok, pass")
		errorz.ValidateTof(&err, false, "error/foo")
		errorz.ValidateToX(&err, 42, false, "error/42")

		errs := errorz.GetErrors(err)
		assert(t, len(errs) == 2).Errorf("❌ expect len(errs) == 2")

		t.Run("format", func(t *testing.T) {
			str := fmt.Sprintf("%#v", err)
			assertEqual(t, str, `
2 errors occurred:
	* error/foo
	* error/42
ezpkg.io/errorz_test/multierr_test.go:███ · TestErrors.func█
testing/testing.go:████ · tRunner
`)
		})
	})
	t.Run("wrap many", func(t *testing.T) {
		errs0 := errorz.Append(
			errors.New("one"),
			errors.New("two"),
			errors.New("three"))
		wErr1 := errorz.Wrap(errs0, "first wrap")
		wErr0 := errorz.Wrapf(wErr1, "outer wrap")

		t.Run("Unwrap", func(t *testing.T) {
			x, ok := errs0.(interface{ Unwrap() []error })
			assert(t, ok).Errorf("❌ expect ok")
			errs := x.Unwrap()
			assert(t, len(errs) == 3).Errorf("❌ expect len(errs) == 3")
			assert(t, fmt.Sprint(errs[0]) == "one").Errorf("❌ expect errs[0] == one")
			assert(t, fmt.Sprint(errs[1]) == "two").Errorf("❌ expect errs[0] == two")
			assert(t, fmt.Sprint(errs[2]) == "three").Errorf("❌ expect errs[0] == three")
		})
		t.Run("GetErrors (nil)", func(t *testing.T) {
			err1 := errorz.GetErrors(nil)
			assert(t, err1 == nil).Errorf("❌ expect err1 == nil")

			err2 := errorz.GetErrors(fmt.Errorf("wrap: %w", errors.New("one")))
			assert(t, err2 == nil).Errorf("❌ expect err2 == nil")
		})
		t.Run("GetErrors", func(t *testing.T) {
			errs := errorz.GetErrors(wErr0)
			fmt.Println(errs)
			assert(t, len(errs) == 3).Errorf("❌ expect len(errs) == 3")
		})
		t.Run("format", func(t *testing.T) {
			t.Run("err0", func(t *testing.T) {
				str := fmt.Sprintf("%#v", errs0)
				fmt.Println(str)
			})
			t.Run("plus", func(t *testing.T) {
				str := fmt.Sprintf("%+v", wErr0)
				fmt.Println(str)
			})
			t.Run("alt", func(t *testing.T) {
				str := fmt.Sprintf("%#v", wErr0)
				fmt.Println(str)
			})
			t.Run("%q", func(t *testing.T) {
				str := fmt.Sprintf("%q", wErr0)
				fmt.Println(str)
			})
		})
	})
}
