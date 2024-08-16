package codez

import (
	"go/ast"
	"log/slog"
	"os"
	"testing"

	. "github.com/onsi/gomega"

	. "ezpkg.io/conveyz"
	"ezpkg.io/logz"
)

var logLevel = slog.LevelDebug

func TestParse(t *testing.T) {
	Ω := GomegaExpect
	log := initLog()
	Convey("parse", t, func() {
		Convey("one", func() {
			stmt, err := parseStmt(log, "a := 42")
			Ω(err).ToNot(HaveOccurred())
			printAst("stmt", nil, stmt)
		})
		Convey("decl", func() {
			decl, err := parseDecl(log, "var a = 42")
			Ω(err).ToNot(HaveOccurred())
			printAst("decl", nil, decl)
		})
		Convey("parseSearch", func() {
			Convey("empty", func() {
				out, err := parseSearch(log, "")
				Ω(err).ToNot(HaveOccurred())
				Ω(out.IsEmpty()).To(BeTrue())
			})
			Convey("comment", func() {
				out, err := parseSearch(log, "// hello")
				Ω(err).ToNot(HaveOccurred())
				Ω(out.IsEmpty()).To(BeTrue())
			})
			Convey("invalid", func() {
				_, err := parseSearch(log, "#")
				Ω(err.Error()).To(ContainSubstring("illegal character"))
			})
			Convey("decl", func() {
				out, err := parseSearch(log, "var a int = 42")
				Ω(err).ToNot(HaveOccurred())

				stmt, _ := out.stmt.(*ast.DeclStmt)
				Ω(stmt).ToNot(BeNil())

				Ω(out.decl).ToNot(BeNil())
			})
			Convey("ident", func() {
				out, err := parseSearch(log, "error")
				Ω(err).ToNot(HaveOccurred())

				ident, _ := out.expr.(*ast.Ident)
				Ω(ident).ToNot(BeNil())
				Ω(out.ident).ToNot(BeNil())
			})
			Convey("file", func() {
				out, err := parseSearch(log, "package main")
				Ω(err).ToNot(HaveOccurred())

				Ω(out.file).ToNot(BeNil())
			})
		})
	})
}

func initLog() logz.Logger {
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	})
	log := slog.New(h)
	return logz.WithEnabler(func(level logz.Level) bool {
		return level >= logLevel
	}).FromLoggerI(log)
}
