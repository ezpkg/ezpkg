package bytez

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type Bytes []byte

func (b Bytes) Bytes() []byte  { return b }
func (b Bytes) String() string { return string(b) }
func (b Bytes) Len() int       { return len(b) }
func (b Bytes) Cap() int       { return cap(b) }
func (b Bytes) Available() int { return cap(b) - len(b) }
func (b *Bytes) Reset()        { *b = (*b)[:0] }
func (b *Bytes) Grow(n int) {
	if n < 0 {
		panic("bytez.Bytes.Grow: negative count")
	}
	if n == 0 {
		return
	}
	L := len(*b)
	*b = append(*b, make([]byte, n)...)
	*b = (*b)[:L]
}
func (b *Bytes) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}
func (b *Bytes) WriteString(s string) (n int, err error) {
	*b = append(*b, []byte(s)...)
	return len(s), nil
}
func (b Bytes) WriteTo(w io.Writer) (n int64, err error) {
	n0, err := w.Write(b)
	return int64(n0), err
}
func (b *Bytes) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}
func (b *Bytes) WriteRune(r rune) (n int, err error) {
	L := len(*b)
	*b = utf8.AppendRune(*b, r)
	return len(*b) - L, nil
}
func (b *Bytes) WriteZ(p []byte) int {
	*b = append(*b, p...)
	return len(p)
}
func (b *Bytes) WriteByteZ(c byte) {
	*b = append(*b, c)
}
func (b *Bytes) WriteRuneZ(r rune) int {
	L := len(*b)
	*b = utf8.AppendRune(*b, r)
	return len(*b) - L
}
func (b *Bytes) WriteStringZ(s string) int {
	*b = append(*b, []byte(s)...)
	return len(s)
}
func (b *Bytes) PrintBytes(p []byte) {
	*b = append(*b, p...)
}
func (b *Bytes) Print(args ...any) {
	_, _ = fmt.Print(args...)
}
func (b *Bytes) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(b, format, args...)
}
func (b *Bytes) Println(args ...any) {
	_, _ = fmt.Fprintln(b, args...)
}
