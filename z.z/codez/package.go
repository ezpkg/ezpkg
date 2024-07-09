package codez

import (
	"golang.org/x/tools/go/packages"
)

type Packages []*Package

type Package struct {
	*packages.Package
}
