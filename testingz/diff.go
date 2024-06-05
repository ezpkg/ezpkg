package testingz

import (
	"diffz"
)

func DiffByChar(expect, actual string) (formatted string, isDiff bool) {
	diffs := diffz.ByChar(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByCharX(expect, actual string, opt diffz.Option) (formatted string, isDiff bool) {
	diffs := diffz.ByCharX(actual, expect, opt)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByCharZ(expect, actual string) (formatted string, isDiff bool) {
	diffs := diffz.ByCharZ(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLine(expect, actual string) (formatted string, isDiff bool) {
	diffs := diffz.ByLine(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLineX(expect, actual string, opt diffz.Option) (formatted string, isDiff bool) {
	diffs := diffz.ByLineX(actual, expect, opt)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLineZ(expect, actual string) (formatted string, isDiff bool) {
	diffs := diffz.ByLineZ(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}
