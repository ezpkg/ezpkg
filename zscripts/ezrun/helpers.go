package main

import (
	"regexp"

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
}
