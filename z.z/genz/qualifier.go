package genz

import (
	"go/types"
)

type DefaultQualifier struct{}

func (q DefaultQualifier) Qualify(pkg *types.Package) string {
	alias := pkg.Name()
	return alias
}
