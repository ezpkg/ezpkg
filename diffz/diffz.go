// Package diffz provides functions for comparing and displaying differences between two strings. It's based on [kylelemons/godebug] and [sergi/go-diff].
//
// [kylelemons/godebug]: https://pkg.go.dev/github.com/kylelemons/godebug
// [sergi/go-diff]: https://pkg.go.dev/github.com/sergi/go-diff
package diffz // import "ezpkg.io/diffz"

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	godebugdiff "github.com/kylelemons/godebug/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type ( // re-export
	Diff      = diffmatchpatch.Diff
	Operation = diffmatchpatch.Operation
)

const ( // re-export
	DiffDelete = diffmatchpatch.DiffDelete
	DiffInsert = diffmatchpatch.DiffInsert
	DiffEqual  = diffmatchpatch.DiffEqual
)

const DefaultPlaceholder = '█'
const DefaultRegexpOpen = "【"
const DefaultRegexpClose = "】"

var diffCharOpts = diffmatchpatch.New()

type Option struct {
	IgnoreSpace bool
	Placeholder rune
	RegexpOpen  string
	RegexpClose string
}

func Default() Option {
	return Option{}
}
func IgnoreSpace() Option {
	return Option{IgnoreSpace: true}
}
func Placeholder() Option {
	return Option{}.AndPlaceholder()
}
func PlaceholderX(placeholder rune) Option {
	return Option{Placeholder: placeholder}
}
func PlaceholderZ() Option {
	return Option{}.AndIgnoreSpace().AndPlaceholder()
}
func (opt Option) AndIgnoreSpace() Option {
	opt.IgnoreSpace = true
	return opt
}
func (opt Option) AndPlaceholder() Option {
	opt.Placeholder = DefaultPlaceholder
	opt.RegexpOpen = DefaultRegexpOpen
	opt.RegexpClose = DefaultRegexpClose
	return opt
}
func (opt Option) AndPlaceholderX(placeholder rune, args ...string) Option {
	opt.Placeholder = placeholder
	switch len(args) {
	case 0:
	case 1:
		panic("expect 2 args for regexp open and close")
	case 2:
		opt.RegexpOpen, opt.RegexpClose = args[0], args[1]
	}
	return opt
}
func (opt Option) DiffByChar(left, right string) (out Diffs) {
	return ByCharX(left, right, opt)
}
func (opt Option) DiffByLine(left, right string) (out Diffs) {
	return ByLineX(left, right, opt)
}

type Diffs struct {
	Items []Diff
}

func (ds Diffs) IsDiff() bool {
	for _, d := range ds.Items {
		if d.Type != diffmatchpatch.DiffEqual {
			return true
		}
	}
	return false
}

func (ds Diffs) String() string {
	return Format(ds)
}

func ByChar(left, right string) Diffs {
	return ByCharX(left, right, Default())
}

func ByCharX(left, right string, opt Option) (out Diffs) {
	out.Items = diffCharOpts.DiffMain(left, right, false)
	return process(opt, out)
}

func ByCharZ(left, right string) Diffs {
	return ByCharX(left, right, PlaceholderZ())
}

func ByLine(left, right string) (out Diffs) {
	return ByLineX(left, right, Default())
}

func ByLineX(left, right string, opt Option) (out Diffs) {
	split := func(s string) (lines []string) {
		if s == "" {
			return nil
		}
		lastIdx := 0
		for i := 0; i < len(s); i++ {
			if s[i] == '\n' {
				lines = append(lines, s[lastIdx:i+1])
				lastIdx = i + 1
			}
		}
		if lastIdx < len(s) {
			lines = append(lines, s[lastIdx:])
		}
		if !strings.HasSuffix(s, "\n") {
			lines[len(lines)-1] += "\n"
		}
		return lines
	}

	lefts, rights := split(left), split(right)
	chunks := godebugdiff.DiffChunks(lefts, rights)
	for _, chunk := range chunks {
		for _, text := range chunk.Deleted {
			out.Items = append(out.Items, Diff{Type: DiffDelete, Text: text})
		}
		for _, text := range chunk.Added {
			out.Items = append(out.Items, Diff{Type: DiffInsert, Text: text})
		}
		for _, text := range chunk.Equal {
			out.Items = append(out.Items, Diff{Type: DiffEqual, Text: text})
		}
	}
	return process(opt, out)
}

func ByLineZ(left, right string) (out Diffs) {
	return ByLineX(left, right, PlaceholderZ())
}

func Format(diffs Diffs) string {
	return diffCharOpts.DiffPrettyText(diffs.Items)
}

func process(opt Option, ds Diffs) Diffs {
	type Error struct {
		left, right string
	}

	if opt == (Option{}) {
		return ds
	}
	if (opt.RegexpOpen == "") != (opt.RegexpClose == "") {
		panic("expect both regexp open and close")
	}

	findRegexp := func(text string) (*regexp.Regexp, string, error) {
		idx := strings.Index(text, opt.RegexpClose)
		if idx < 0 {
			return nil, "", errors.New("《REGEXP:NO_CLOSE》")
		}
		expr := strings.TrimSpace(text[:idx])
		if expr == "" {
			return nil, "", errors.New("《REGEXP:EMPTY》")
		}
		re, err := regexp.Compile(expr)
		if err != nil {
			return nil, text[:idx], errors.New("《REGEXP:INVALID》")
		}
		return re, text[:idx], nil
	}
	processRegexpDiffs := func(xDiffs []Diff, op Operation) (opIdx, negIdx int) {
		for i, diff := range xDiffs {
			switch diff.Type {
			case DiffEqual:
				opIdx += len(diff.Text)
				negIdx += len(diff.Text)
			case op:
				hasOpen := strings.Contains(diff.Text, opt.RegexpOpen)
				hasClose := strings.Contains(diff.Text, opt.RegexpClose)
				switch {
				case hasOpen && !hasClose:
					opIdx += len(diff.Text)
				case hasClose:
					opIdx += len(diff.Text)
					if i+1 < len(xDiffs) && xDiffs[i+1].Type == -op {
						negIdx += len(xDiffs[i+1].Text)
					}
					return opIdx, negIdx
				}
			case -op:
				negIdx += len(diff.Text)
			}
		}
		return opIdx, negIdx
	}
	match := func(delText, insText string) (delNL, insNL int, ok bool, _ *Error) {
		delI, insI := 0, 0
		delL, insL := len(delText), len(insText)
		for {
			delCh, delSize := utf8.DecodeRuneInString(delText[delI:])
			if delSize == 0 {
				return delI, insI, delI == delL || insI == insL, nil
			}
			insCh, insSize := utf8.DecodeRuneInString(insText[insI:])
			if insSize == 0 {
				return delI, insI, delI == delL || insI == insL, nil
			}
			if delCh == utf8.RuneError {
				return 0, 0, false, &Error{left: "《RUNE ERROR》"}
			}
			if insCh == utf8.RuneError {
				return 0, 0, false, &Error{right: "《RUNE ERROR》"}
			}
			if delCh == '\n' && insCh == '\n' {
				delNL, insNL = delI+1, insI+1
			}
			if delCh == insCh {
				delI, insI = delI+delSize, insI+insSize
				continue
			}
			if opt.IgnoreSpace && isSpace(delCh) {
				delI += delSize
				continue
			}
			if opt.IgnoreSpace && isSpace(insCh) {
				insI += insSize
				continue
			}
			if opt.Placeholder != 0 {
				if opt.Placeholder == delCh || opt.Placeholder == insCh {
					delI += delSize
					insI += insSize
					continue
				}
			}
			if opt.RegexpOpen != "" && opt.RegexpClose != "" {
				if !strings.Contains(delText[delI:], opt.RegexpOpen) && !strings.Contains(insText[insI:], opt.RegexpOpen) {
					return delNL, insNL, false, nil
				}

				xDiffs := diffCharOpts.DiffMain(delText[delI:], insText[insI:], false)
				i := 0
				for ; i < len(xDiffs) && xDiffs[i].Type == DiffEqual; i++ {
					delI += len(xDiffs[i].Text)
					insI += len(xDiffs[i].Text)
				}
				xDiffs = xDiffs[i:]

				switch {
				case len(xDiffs) == 0:
					continue

				case strings.HasPrefix(delText[delI:], opt.RegexpOpen):
					delIdx, insIdx := processRegexpDiffs(xDiffs, DiffDelete)
					delDiff := delText[delI : delI+delIdx]
					insDiff := insText[insI : insI+insIdx]
					re, rawRe, err := findRegexp(delDiff[len(opt.RegexpOpen):])
					if err != nil {
						return delNL, insNL, false, &Error{right: err.Error()}
					}
					loc := re.FindStringIndex(insDiff)
					if loc == nil || loc[0] != 0 {
						return delNL, insNL, false, nil
					}
					delI += len(opt.RegexpOpen) + len(rawRe) + len(opt.RegexpClose)
					insI += loc[1] - loc[0]
					continue

				case strings.HasPrefix(insText[insI:], opt.RegexpOpen):
					insIdx, delIdx := processRegexpDiffs(xDiffs, DiffInsert)
					insDiff := insText[insI : insI+insIdx]
					delDiff := delText[delI : delI+delIdx]
					re, rawRe, err := findRegexp(insDiff[len(opt.RegexpOpen):])
					if err != nil {
						return delNL, insNL, false, &Error{left: err.Error()}
					}
					loc := re.FindStringIndex(delDiff)
					if loc == nil || loc[0] != 0 {
						return delNL, insNL, false, nil
					}
					insI += len(opt.RegexpOpen) + len(rawRe) + len(opt.RegexpClose)
					delI += loc[1] - loc[0]
					continue
				}
			}
			return delNL, insNL, false, nil
		}
	}

	i, iDel, iIns, L := 0, 0, 0, len(ds.Items)
	outItems := make([]Diff, 0, len(ds.Items))
	remainDel, remainIns := "", ""

	appendOut := func(op Operation, text string) {
		if opt.IgnoreSpace && strings.TrimSpace(text) == "" {
			return
		}
		if text != "" {
			outItems = append(outItems, Diff{Type: op, Text: text})
		}
	}
	appendEqual := func(delText, insText string) {
		if delText != "" {
			outItems = append(outItems, Diff{Type: DiffEqual, Text: delText})
		}
		// TODO: handle placeholder when format
	}
	skipEqualOrSpaceDiffs := func() {
		for ; i < L; i++ {
			diff := ds.Items[i]
			if diff.Type == DiffEqual {
				outItems = append(outItems, diff)
				continue
			}
			if opt.IgnoreSpace && strings.TrimSpace(diff.Text) == "" {
				continue
			}
			return
		}
	}
	nextDelDiff := func() (out string, ok bool) {
		if remainDel != "" {
			out, remainDel = remainDel, ""
			return out, true
		}
		iDel = max(iDel, i)
		for ; iDel < L; iDel++ {
			diff := ds.Items[iDel]
			if opt.IgnoreSpace && strings.TrimSpace(diff.Text) == "" {
				continue
			}
			switch diff.Type {
			case DiffEqual:
				if opt.IgnoreSpace && strings.TrimSpace(diff.Text) == "" {
					continue
				}
				return "", false
			case DiffDelete:
				iDel++
				return diff.Text, true
			}
		}
		return "", false
	}
	nextInsDiff := func() (out string, ok bool) {
		if remainIns != "" {
			out, remainIns = remainIns, ""
			return out, true
		}
		iIns = max(iIns, i)
		for ; iIns < L; iIns++ {
			diff := ds.Items[iIns]
			if opt.IgnoreSpace && strings.TrimSpace(diff.Text) == "" {
				continue
			}
			switch diff.Type {
			case DiffEqual:
				return "", false
			case DiffInsert:
				iIns++
				return diff.Text, true
			}
		}
		return "", false
	}
	processRemaining := func() {
		appendOut(DiffDelete, remainDel)
		appendOut(DiffInsert, remainIns)
		remainDel, remainIns = "", ""
		i = max(i, min(iDel, iIns))
		for ; i < L; i++ {
			diff := ds.Items[i]
			switch diff.Type {
			case DiffEqual:
				return
			case DiffDelete:
				if i >= iDel {
					iDel++
					appendOut(diff.Type, diff.Text)
				}
			case DiffInsert:
				if i >= iIns {
					iIns++
					appendOut(diff.Type, diff.Text)
				}
			}
		}
	}

	for i < L {
		skipEqualOrSpaceDiffs()
		for {
			delText, ok := nextDelDiff()
			if !ok {
				processRemaining()
				break
			}
			insText, ok := nextInsDiff()
			if !ok {
				remainDel = delText
				processRemaining()
				break
			}
			delNL, insNL, ok, err := match(delText, insText)
			appendEqual(delText[:delNL], insText[:insNL])
			remainDel, remainIns = delText[delNL:], insText[insNL:]
			switch {
			case err != nil && err.left != "":
				appendOut(DiffDelete, "❌"+err.left)
				processRemaining()
				break
			case err != nil && err.right != "":
				appendOut(DiffInsert, "❌"+err.right)
				processRemaining()
				break
			case !ok:
				processRemaining()
				break
			}
		}
	}
	return Diffs{Items: outItems}
}

func isSpace(c rune) bool {
	return unicode.IsSpace(c)
}
