package ggen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDirectivesFromBody(t *testing.T) {
	t.Run("go:build directive", func(t *testing.T) {
		body := `//go:build tag1,tag2

package main
`
		var directives, inlineDirectives []Directive
		errs := parseDirectivesFromBody([]byte(body), &directives, &inlineDirectives)

		require.Len(t, errs, 0)
		require.Len(t, directives, 1)
		require.Len(t, inlineDirectives, 0)
		require.Equal(t, Directive{
			Raw: "go:build tag1,tag2",
			Cmd: "go:build",
			Arg: "tag1,tag2",
		}, directives[0])
	})

	t.Run("+sample directive", func(t *testing.T) {
		body := `
package main

//+sample
`
		var directives, inlineDirectives []Directive
		errs := parseDirectivesFromBody([]byte(body), &directives, &inlineDirectives)

		require.Len(t, errs, 0)
		require.Len(t, directives, 1)
		require.Len(t, inlineDirectives, 0)
		require.Equal(t, Directive{
			Raw: "+sample",
			Cmd: "sample",
			Arg: "",
		}, directives[0])
	})
}
