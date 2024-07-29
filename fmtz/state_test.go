package fmtz_test

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"ezpkg.io/bytez"
	"ezpkg.io/fmtz"
	"ezpkg.io/stringz"
)

// ASSERT: All packages fmtz,bytez,stringz implement the same interface.
var _ WriterZ = fmtz.State{}
var _ WriterZ = fmtz.MustState{}
var _ WriterZ = &bytez.Buffer{}
var _ WriterZ = &bytez.Bytes{}
var _ WriterZ = &stringz.Builder{}

type WriterZ interface {
	Write(p []byte) (int, error)
	WriteByte(c byte) error
	WriteRune(r rune) (int, error)
	WriteString(s string) (int, error)
	WriteZ(p []byte) int
	WriteByteZ(c byte)
	WriteRuneZ(r rune) int
	WriteStringZ(s string) int
	PrintBytes(p []byte)
	Print(args ...any)
	Printf(format string, args ...any)
	Println(args ...any)
}

type CodeStd Code
type Code struct {
	Char   rune
	Number int
}

func (c Code) Format(s0 fmt.State, r rune) {
	s := fmtz.WrapState(s0)
	s.WriteRuneZ(c.Char)
	s.Print(c.Number)
}
func (c CodeStd) Format(s fmt.State, r rune) {
	var p []byte
	p = utf8.AppendRune(p, c.Char)
	_, _ = s.Write(p)
	_, _ = fmt.Fprint(s, c.Number)
}

func TestState(t *testing.T) {
	code := Code{'Ω', 123}
	out := fmt.Sprint(code)
	if out != "Ω123" {
		t.Errorf("unexpected output: %q", out)
	}
}

func BenchmarkState(b *testing.B) {
	code := Code{'Ω', 123}
	codeStd := CodeStd(code)

	var out string
	b.Run("std", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			out = fmt.Sprint(codeStd)
		}
		b.StopTimer()
		if out != "Ω123" {
			b.Fatalf("unexpected output: %s", out)
		}
	})
	b.Run("fmtz", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			out = fmt.Sprint(code)
		}
		b.StopTimer()
		if out != "Ω123" {
			b.Fatalf("unexpected output: %q", out)
		}
	})
}
