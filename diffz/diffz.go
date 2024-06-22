// Package diffz provides functions for comparing and displaying differences between two strings. It's based on kylelemons/godebug and sergi/go-diff.
package diffz // import "ezpkg.io/diffz"

import (
	"strings"
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

type Option struct {
	IgnoreSpace bool
	Placeholder rune
}

func Default() Option {
	return Option{}
}
func IgnoreSpace() Option {
	return Option{IgnoreSpace: true}
}
func Placeholder() Option {
	return Option{Placeholder: DefaultPlaceholder}
}
func PlaceholderX(placeholder rune) Option {
	return Option{Placeholder: placeholder}
}
func PlaceholderZ() Option {
	return Option{
		IgnoreSpace: true,
		Placeholder: DefaultPlaceholder,
	}
}
func (opt Option) AndIgnoreSpace() Option {
	opt.IgnoreSpace = true
	return opt
}
func (opt Option) AndPlaceholder() Option {
	opt.Placeholder = DefaultPlaceholder
	return opt
}
func (opt Option) AndPlaceholderX(placeholder rune) Option {
	opt.Placeholder = placeholder
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

func ByChar(left, right string) Diffs {
	return ByCharX(left, right, Default())
}

func ByCharX(left, right string, opt Option) (out Diffs) {
	dmp := diffmatchpatch.New()
	out.Items = dmp.DiffMain(left, right, false)
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
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyText(diffs.Items)
}

func process(opt Option, ds Diffs) Diffs {
	if opt == (Option{}) {
		return ds
	}
	match := func(delText, insText string) (delNL, insNL int, ok bool) {
		delI, insI := 0, 0
		delL, insL := len(delText), len(insText)
		for {
			delCh, delSize := utf8.DecodeRuneInString(delText[delI:])
			if delSize == 0 {
				return delI, insI, delI == delL || insI == insL
			}
			insCh, insSize := utf8.DecodeRuneInString(insText[insI:])
			if insSize == 0 {
				return delI, insI, delI == delL || insI == insL
			}
			if delCh == utf8.RuneError || insCh == utf8.RuneError {
				return 0, 0, false
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
			if opt.Placeholder != 0 && opt.Placeholder == delCh || opt.Placeholder == insCh {
				delI += delSize
				insI += insSize
				continue
			}
			return delNL, insNL, false
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
			delNL, insNL, ok := match(delText, insText)
			appendEqual(delText[:delNL], insText[:insNL])
			remainDel, remainIns = delText[delNL:], insText[insNL:]
			if !ok {
				processRemaining()
				break
			}
		}
	}
	return Diffs{Items: outItems}
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
