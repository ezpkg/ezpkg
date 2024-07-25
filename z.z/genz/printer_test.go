package ggen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	p := &printer{
		pkgPathByAlias: make(map[string]string),
		aliasByPkgPath: make(map[string]string),
	}
	p.Import("one", "encoding/one")
	p.Import("one", "example.com/one")
	p.Import("one", "github.com/one/one")

	assert.Equal(t, "encoding/one", p.pkgPathByAlias["one"])
	assert.Equal(t, "example.com/one", p.pkgPathByAlias["one1"])
	assert.Equal(t, "github.com/one/one", p.pkgPathByAlias["one2"])

	assert.Equal(t, "one", p.aliasByPkgPath["encoding/one"])
	assert.Equal(t, "one1", p.aliasByPkgPath["example.com/one"])
	assert.Equal(t, "one2", p.aliasByPkgPath["github.com/one/one"])
}
