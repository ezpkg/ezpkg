package codez

import (
	"go/ast"
	"log/slog"
	"os"
	"testing"

	g "github.com/onsi/gomega"

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
			Ω(err).ToNot(g.HaveOccurred())
			printAst("stmt", nil, stmt)
		})
		Convey("decl", func() {
			decl, err := parseDecl(log, "var a = 42")
			Ω(err).ToNot(g.HaveOccurred())
			printAst("decl", nil, decl)
		})
		Convey("parseSearch", func() {
			Convey("empty", func() {
				out, err := parseSearch(log, "")
				Ω(err).ToNot(g.HaveOccurred())
				Ω(out.IsEmpty()).To(g.BeTrue())
			})
			Convey("comment", func() {
				out, err := parseSearch(log, "// hello")
				Ω(err).ToNot(g.HaveOccurred())
				Ω(out.IsEmpty()).To(g.BeTrue())
			})
			Convey("invalid", func() {
				_, err := parseSearch(log, "#")
				Ω(err.Error()).To(g.ContainSubstring("illegal character"))
			})
			Convey("decl", func() {
				out, err := parseSearch(log, "var a int = 42")
				Ω(err).ToNot(g.HaveOccurred())

				stmt, _ := out.stmt.(*ast.DeclStmt)
				Ω(stmt).ToNot(g.BeNil())

				Ω(out.decl).ToNot(g.BeNil())
			})
			Convey("ident", func() {
				out, err := parseSearch(log, "error")
				Ω(err).ToNot(g.HaveOccurred())

				ident, _ := out.expr.(*ast.Ident)
				Ω(ident).ToNot(g.BeNil())
				Ω(out.ident).ToNot(g.BeNil())
			})
			Convey("file", func() {
				out, err := parseSearch(log, "package main")
				Ω(err).ToNot(g.HaveOccurred())

				Ω(out.file).ToNot(g.BeNil())
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
