package codez

import (
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
	})
}

func initLog() logz.Logger {
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	})
	log := slog.New(h)
	return logz.WithEnabler(func(level logz.Level) bool {
		return level.ToInt() >= int(logLevel)
	}).FromLoggerI(log)
}
