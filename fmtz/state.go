package fmtz

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type State struct {
	fmt.State
}

func WrapState(state fmt.State) State {
	return State{State: state}
}
func (b State) WriteByte(c byte) error {
	d := [1]byte{c}
	_, err := b.Write(d[:])
	return err
}
func (b State) WriteRune(r rune) (int, error) {
	d := [utf8.UTFMax]byte{}
	return b.Write(utf8.AppendRune(d[:0], r))
}
func (b State) WriteString(s string) (int, error) {
	return io.WriteString(b.State, s)
}
func (b State) WriteZ(p []byte) int {
	n, _ := b.State.Write(p)
	return n
}
func (b State) WriteByteZ(c byte) {
	_ = b.WriteByte(c)
}
func (b State) WriteRuneZ(r rune) int {
	n, _ := b.WriteRune(r)
	return n
}
func (b State) WriteStringZ(s string) int {
	n, _ := io.WriteString(b.State, s)
	return n
}
func (b State) PrintBytes(p []byte) {
	_, _ = b.Write(p)
}
func (b State) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(b.State, format, args...)
}
func (b State) Println(args ...any) {
	_, _ = fmt.Fprintln(b.State, args...)
}
func (b State) Print(args ...any) {
	_, _ = fmt.Fprint(b.State, args...)
}

type MustState struct {
	fmt.State
}

func WrapMustState(state fmt.State) MustState {
	return MustState{State: state}
}
func (b MustState) WriteByte(c byte) error {
	d := [1]byte{c}
	_, err := b.Write(d[:])
	return err
}
func (b MustState) WriteRune(r rune) (int, error) {
	d := [utf8.UTFMax]byte{}
	return b.Write(utf8.AppendRune(d[:], r))
}
func (b MustState) WriteString(s string) (int, error) {
	return io.WriteString(b.State, s)
}
func (b MustState) WriteZ(p []byte) int {
	return must(b.State.Write(p))
}
func (b MustState) WriteByteZ(c byte) {
	mustZ(b.WriteByte(c))
}
func (b MustState) WriteRuneZ(r rune) int {
	return must(b.WriteRune(r))
}
func (b MustState) WriteStringZ(s string) int {
	return must(io.WriteString(b.State, s))
}
func (b MustState) PrintBytes(p []byte) {
	must(b.Write(p))
}
func (b MustState) Printf(format string, args ...any) {
	must(fmt.Fprintf(b.State, format, args...))
}
func (b MustState) Println(args ...any) {
	must(fmt.Fprintln(b.State, args...))
}
func (b MustState) Print(args ...any) {
	must(fmt.Fprint(b.State, args...))
}
