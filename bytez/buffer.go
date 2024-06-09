package bytez

import (
	"bytes"
	"fmt"
	"io"
)

type Buffer bytes.Buffer

func (b *Buffer) unwrap() *bytes.Buffer                          { return (*bytes.Buffer)(b) }
func (b *Buffer) Bytes() []byte                                  { return b.unwrap().Bytes() }
func (b *Buffer) AvailableBuffer() []byte                        { return b.unwrap().AvailableBuffer() }
func (b *Buffer) String() string                                 { return b.unwrap().String() }
func (b *Buffer) Len() int                                       { return b.unwrap().Len() }
func (b *Buffer) Cap() int                                       { return b.unwrap().Cap() }
func (b *Buffer) Available() int                                 { return b.unwrap().Available() }
func (b *Buffer) Truncate(n int)                                 { b.unwrap().Truncate(n) }
func (b *Buffer) Reset()                                         { b.unwrap().Reset() }
func (b *Buffer) Grow(n int)                                     { b.unwrap().Grow(n) }
func (b *Buffer) Write(p []byte) (n int, err error)              { return b.unwrap().Write(p) }
func (b *Buffer) WriteString(s string) (n int, err error)        { return b.unwrap().WriteString(s) }
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error)      { return b.unwrap().ReadFrom(r) }
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error)       { return b.unwrap().WriteTo(w) }
func (b *Buffer) WriteByte(c byte) error                         { return b.unwrap().WriteByte(c) }
func (b *Buffer) WriteRune(r rune) (n int, err error)            { return b.unwrap().WriteRune(r) }
func (b *Buffer) Read(p []byte) (n int, err error)               { return b.unwrap().Read(p) }
func (b *Buffer) Next(n int) []byte                              { return b.unwrap().Next(n) }
func (b *Buffer) ReadByte() (byte, error)                        { return b.unwrap().ReadByte() }
func (b *Buffer) ReadRune() (r rune, size int, err error)        { return b.unwrap().ReadRune() }
func (b *Buffer) UnreadRune() error                              { return b.unwrap().UnreadRune() }
func (b *Buffer) UnreadByte() error                              { return b.unwrap().UnreadByte() }
func (b *Buffer) ReadBytes(delim byte) (line []byte, err error)  { return b.unwrap().ReadBytes(delim) }
func (b *Buffer) ReadString(delim byte) (line string, err error) { return b.unwrap().ReadString(delim) }

func (b *Buffer) WriteZ(p []byte) int {
	n, _ := b.unwrap().Write(p)
	return n
}
func (b *Buffer) WriteByteZ(c byte) {
	b.unwrap().WriteByte(c)
}
func (b *Buffer) WriteRuneZ(r rune) int {
	n, _ := b.unwrap().WriteRune(r)
	return n
}
func (b *Buffer) WriteStringZ(s string) {
	_, _ = b.unwrap().WriteString(s)
}
func (b *Buffer) PrintBytes(p []byte) {
	_, _ = b.Write(p)
}
func (b *Buffer) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(b.unwrap(), format, args...)
}
func (b *Buffer) Println(args ...any) {
	_, _ = fmt.Fprintln(b.unwrap(), args...)
}
func (b *Buffer) Print(args ...any) {
	_, _ = fmt.Fprint(b.unwrap(), args...)
}
