package diffz

import (
	"regexp"
	"strings"

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

const DefaultPlaceholder = "█"

var regexpSpace = regexp.MustCompile(`^\s*$`)

type Diffs struct {
	Items []diffmatchpatch.Diff

	ignoreSpace bool
	placeholder string
}

func (ds Diffs) Unwrap() []diffmatchpatch.Diff {
	return ds.Items
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
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(left, right, false)
	return Diffs{Items: diffs}
}

func ByCharZ(left, right string) Diffs {
	return ByChar(left, right).IgnoreSpace().Placeholder()
}

func ByLine(left, right string) (out Diffs) {
	chunks := godebugdiff.DiffChunks(split(left), split(right))
	for _, chunk := range chunks {
		for _, line := range chunk.Added {
			out.Items = append(out.Items, Diff{
				Type: DiffInsert,
				Text: line,
			})
		}
		for _, line := range chunk.Deleted {
			out.Items = append(out.Items, Diff{
				Type: DiffDelete,
				Text: line,
			})
		}
		for _, line := range chunk.Equal {
			out.Items = append(out.Items, Diff{
				Type: DiffEqual,
				Text: line,
			})
		}
	}
	return out
}

func Format(diffs Diffs) string {
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyText(diffs.Items)
}

func (ds Diffs) IgnoreSpace() Diffs {
	ds.ignoreSpace = true
	outItems := make([]diffmatchpatch.Diff, 0, len(ds.Items))
	for i := 0; i < len(ds.Items); i++ {
		diff := ds.Items[i]
		switch diff.Type {
		case DiffEqual:
			outItems = append(outItems, diff)

		case DiffInsert:
			if strings.TrimSpace(diff.Text) == "" {
				continue // ignore empty diff
			}
			outItems = append(outItems, diff)

		case DiffDelete:
			if strings.TrimSpace(diff.Text) == "" {
				continue // ignore empty diff
			}
			if len(ds.Items) <= i+1 {
				outItems = append(outItems, diff)
				continue
			}
			nextDiff := ds.Items[i+1]
			if nextDiff.Type != DiffInsert {
				outItems = append(outItems, diff)
				continue
			}
			if ds.matchPlaceholder(diff.Text, nextDiff.Text) {
				i++ // ignore the diff after the placeholder
				outItems = append(outItems, Diff{
					Type: DiffEqual,
					Text: nextDiff.Text,
				})
				continue
			}
			if ds.matchPlaceholder(nextDiff.Text, diff.Text) {
				i++ // ignore the diff after the placeholder
				outItems = append(outItems, Diff{
					Type: DiffEqual,
					Text: diff.Text,
				})
				continue
			}
			outItems = append(outItems, diff)
		}
	}
	ds.Items = outItems
	return ds
}

func (ds Diffs) Placeholder() Diffs {
	return ds.PlaceholderX(DefaultPlaceholder)
}

func (ds Diffs) PlaceholderX(placeholder string) Diffs {
	if placeholder == "" {
		panic("placeholder cannot be empty")
	}
	ds.placeholder = placeholder
	outItems := make([]Diff, 0, len(ds.Items))
	for i := 0; i < len(ds.Items); i++ {
		diff := ds.Items[i]
		switch diff.Type {
		case DiffEqual:
			outItems = append(outItems, diff)

		case DiffInsert:
			outItems = append(outItems, diff)

		case DiffDelete:
			if len(ds.Items) <= i+1 {
				outItems = append(outItems, diff)
				continue
			}
			nextDiff := ds.Items[i+1]
			if nextDiff.Type != DiffInsert {
				outItems = append(outItems, diff)
				continue
			}
			if ds.matchPlaceholder(diff.Text, nextDiff.Text) {
				i++ // ignore the diff after the placeholder
				outItems = append(outItems, Diff{
					Type: DiffEqual,
					Text: nextDiff.Text,
				})
				continue
			}
			if ds.matchPlaceholder(nextDiff.Text, diff.Text) {
				i++ // ignore the diff after the placeholder
				outItems = append(outItems, Diff{
					Type: DiffEqual,
					Text: diff.Text,
				})
				continue
			}
			outItems = append(outItems, diff)

		default:
			panic("unreachable")
		}
	}
	ds.Items = outItems
	return ds
}

func (ds Diffs) matchPlaceholder(placeholderText, normalText string) bool {
	if ds.placeholder == "" {
		return false
	}
	if ds.ignoreSpace {
		placeholderText = strings.TrimSpace(placeholderText)
		normalText = strings.TrimSpace(normalText)
	}
	placeholder := ds.placeholder
	if len(placeholderText)%len(placeholder) != 0 {
		return false
	}
	N := len(placeholderText) / len(placeholder)
	for i := 0; i < len(placeholderText); i += len(placeholder) {
		if placeholderText[i:i+len(placeholder)] != placeholder {
			return false
		}
	}
	count := 0
	for range normalText {
		count++
		if count > N {
			return false
		}
	}
	return count == N
}

func split(s string) (lines []string) {
	lastIdx := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[lastIdx:i+1])
			lastIdx = i + 1
		}
	}
	return lines
}
