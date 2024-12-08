package main

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
)

// ezpkg.json
type EzpkgInfo struct {
	Version string `json:"version"`
}

func (x *EzpkgInfo) Validate() (err error) {
	reVer := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	errorz.ValidateTof(&err, x.Version != "", "version is required")
	errorz.ValidateTof(&err, reVer.MatchString(x.Version), "malformed version %q", x.Version)
	return err
}

type PkgInfo struct {
	*packages.Package
	pkgDir string
	goMod  *modfile.File

	goModPublish []byte
	goModLocal   []byte

	ezDeps    []string // direct ezpkg.io dependencies in format "ezpkg.io/..."
	ezDepsAll []string // direct and indirect ezpkg.io dependencies
}

func splitLine(s string, prefix string) (line, remain string) {
	if !strings.HasPrefix(s, prefix) {
		return "", s
	}
	remain = s[len(prefix):]
	idx := strings.Index(remain, "\n")
	idx = max(idx, 0)
	return s[:len(prefix)+idx], strings.TrimSpace(remain[idx:])
}
func panicf(format string, args ...any) {
	panic(fmt.Errorf(format, args...))
}
