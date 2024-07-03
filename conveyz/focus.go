package conveyz

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"ezpkg.io/typez"
	"ezpkg.io/unsafez"
)

var (
	skippedTests      bool
	cachedFocusConvey int // 0: not init, 1: has focus, -1: no focus

	fileMap = map[string]*testFile{}
)

type testFile struct {
	Path  string
	Lines []testLine
}

type testLine struct {
	Indent  int
	Text    string
	FConvey bool
	Convey  bool
}

// detect if any *_test.go file has FocusConvey (or FConvey)
func pkgHasFocusConvey() bool {
	if cachedFocusConvey != 0 {
		return cachedFocusConvey == 1
	}

	_, currentFile, _, ok := runtime.Caller(2)
	if !ok {
		return false
	}
	if !strings.HasSuffix(currentFile, "_test.go") {
		panic(fmt.Sprintf("UNEXPECTED: %q does not end with _test.go", currentFile))
	}

	detected := parseDir(filepath.Dir(currentFile))
	cachedFocusConvey = typez.If(detected, 1, -1)
	return detected
}

// find FocusConvey in all dirs, parse and store them into fileMap
func parseDir(dir string) (detected bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		if strings.IndexAny(entry.Name(), "_.") == 0 {
			continue // files start with "_" or "." are ignored
		}
		path := filepath.Join(dir, entry.Name())
		data, err0 := os.ReadFile(path)
		if err0 != nil {
			panic(fmt.Sprintf("UNEXPECTED: failed to open file: %v", err0))
		}
		dataStr := unsafez.BytesToString(data)
		file, detected0 := parseFile(path, dataStr)
		fileMap[path] = file
		detected = detected || detected0
	}
	return detected
}

func parseFile(path, data string) (_ *testFile, detected bool) {
	tf := &testFile{Path: path}
	for _, line := range strings.Split(data, "\n") {
		tl := testLine{
			Indent: countTabs(line),
			Text:   line,
		}
		if strings.Contains(line, "FConvey(") || strings.Contains(line, "FocusConvey(") {
			detected = true
			tl.FConvey = true
		} else if strings.Contains(line, "Convey(") {
			tl.Convey = true
		}
		tf.Lines = append(tf.Lines, tl)
	}
	return tf, detected
}

// CASE 1:
// 1A. detect if any child block inside this convey block has FocusConvey (or FConvey)
// 1B. detect if any parent block containing this convey block has FocusConvey (or FConvey)
//
//	Convey        -> convert to FocusConvey (parent)
//	  Convey      -> not
//	    Convey    -> not
//	  FConvey     -> convert to FocusConvey (.)
//	    Convey    -> convert to FocusConvey (child)
//	  Convey      -> not
//	    Convey    -> not
//
// -> return 1: convert the Convey block to FocusConvey
//
// CASE 2: detect if there is no FConvey or FocusConvey in the whole file
// -> return -1: convert the Convey block to SkipConvey
func shouldConvert() (out int) {
	_, file, lineNumber, ok := runtime.Caller(2)
	if !ok {
		return 0
	}
	tf := fileMap[file]
	if tf == nil {
		return 0 // leave it as is
	}

	lineIndex := lineNumber - 1
	{
		// CASE 1: detect FocusConvey in child or parent blocks
		// -> return 1: convert the Convey block to FocusConvey
		line := tf.Lines[lineIndex]
		// check if the child code lines inside this convey block has FocusConvey (or FConvey)
		for i := lineNumber; i < len(tf.Lines); i++ {
			if tf.Lines[i].Indent <= line.Indent && tf.Lines[i].Indent != 0 /* skip empty lines */ {
				break
			}
			if tf.Lines[i].FConvey {
				return 1 // convert to FocusConvey
			}
		}
		// check if the parent code lines of this convey block has FocusConvey (or FConvey)
		for parent := getParent(tf, lineIndex); parent > 0; parent = getParent(tf, parent) {
			if tf.Lines[parent].FConvey {
				return 1 // convert to FocusConvey
			}
		}
	}
	{
		// CASE 2: detect if there is no FConvey or FocusConvey in the whole file
		// -> return -1: convert the Convey block to SkipConvey
		detected := false
		for _, line := range tf.Lines {
			if line.FConvey {
				detected = true
				break
			}
		}
		if !detected {
			return -1 // convert to SkipConvey
		}
	}
	return 0
}

func getParent(tf *testFile, lineIndex int) int {
	line := tf.Lines[lineIndex]
	for i := lineIndex - 1; i >= 0; i-- {
		if tf.Lines[i].Indent == 0 {
			// skip empty lines
			continue
		}
		if tf.Lines[i].Indent >= line.Indent {
			// skip lines with same or greater indent
			continue
		}
		if tf.Lines[i].Indent < line.Indent {
			// found the parent
			return i
		}
	}
	return 0
}

func countTabs(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] != '\t' {
			return i
		}
	}
	return len(s)
}
