package stringz

import (
	"fmt"
	"strings"
)

type Builder strings.Builder

func (b *Builder) unwrap() *strings.Builder          { return (*strings.Builder)(b) }
func (b *Builder) String() string                    { return b.unwrap().String() }
func (b *Builder) Len() int                          { return b.unwrap().Len() }
func (b *Builder) Cap() int                          { return b.unwrap().Cap() }
func (b *Builder) Reset()                            { b.unwrap().Reset() }
func (b *Builder) Grow(n int)                        { b.unwrap().Grow(n) }
func (b *Builder) Write(p []byte) (int, error)       { return b.unwrap().Write(p) }
func (b *Builder) WriteByte(c byte) error            { return b.unwrap().WriteByte(c) }
func (b *Builder) WriteRune(r rune) (int, error)     { return b.unwrap().WriteRune(r) }
func (b *Builder) WriteString(s string) (int, error) { return b.unwrap().WriteString(s) }

func (b *Builder) WriteZ(p []byte) int {
	n, _ := b.unwrap().Write(p)
	return n
}
func (b *Builder) WriteByteZ(c byte) {
	b.unwrap().WriteByte(c)
}
func (b *Builder) WriteRuneZ(r rune) int {
	n, _ := b.unwrap().WriteRune(r)
	return n
}
func (b *Builder) WriteStringZ(s string) int {
	n, _ := b.unwrap().WriteString(s)
	return n
}
func (b *Builder) PrintBytes(p []byte) {
	_, _ = b.Write(p)
}
func (b *Builder) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(b.unwrap(), format, args...)
}
func (b *Builder) Println(args ...any) {
	_, _ = fmt.Fprintln(b.unwrap(), args...)
}
func (b *Builder) Print(args ...any) {
	_, _ = fmt.Fprint(b.unwrap(), args...)
}
