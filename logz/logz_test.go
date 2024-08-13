package logz

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	. "ezpkg.io/conveyz"
	"ezpkg.io/diffz"
	"ezpkg.io/stringz"
)

var (
	_ LoggerP    = (*log.Logger)(nil)
	_ LoggerI    = (*slog.Logger)(nil)
	_ logger0ctx = (*slog.Logger)(nil)
	_ Loggerw    = (*zap.SugaredLogger)(nil)
	_ Loggerf    = (*zap.SugaredLogger)(nil)
)

func Test(t *testing.T) {
	Ω := GomegaExpect
	Convey("Test", t, func() {
		var b stringz.Builder

		Convey("slog", func() {
			lv := &slog.LevelVar{}
			lv.Set(slog.LevelDebug)

			opt := &slog.HandlerOptions{Level: lv}
			handler := slog.NewTextHandler(&b, opt)
			logger := slog.New(handler)
			loggerz := FromLoggerI(logger)

			assert := func() {
				s := b.String()
				printLog(s)
				Ω(strings.HasPrefix(s, "time=")).To(BeTrue())
				Ω(strings.HasSuffix(s, `level=INFO msg="Hello, World!" name=Alice age=18`))
			}
			Convey("slog.Info", func() {
				logger.Info("Hello, World!", "name", "Alice", "age", "18")

				assert()
			})
			Convey("logz.Info", func() {
				loggerz.Infow("Hello, World!", "name", "Alice", "age", 18)

				assert()
			})
			Convey("enabler", func() {
				ctx := context.Background()
				Convey("inner logger implements Enabler", func() {
					Convey("default", func() {
						lv.Set(slog.LevelDebug)
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelWarn)).To(BeTrue())

						lv.Set(slog.LevelInfo)
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelWarn)).To(BeTrue())

						lv.Set(slog.LevelWarn)
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelWarn)).To(BeTrue())
					})
				})
				Convey("inner logger does not implement Enabler", func() {
					simpleLogger := log.New(&b, "", 0)
					loggerz := FromLoggerP(simpleLogger)

					Convey("default", func() {
						// default to info
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeTrue())

						SetDefaultEnableLevel(LevelDebug)
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeTrue())

						SetDefaultEnableLevel(LevelInfo)
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeTrue())
					})
					Convey("option", func() {
						loggerz := WithEnabler(func(level Level) bool {
							return level >= lv.Level()
						}).FromLoggerI(logger)

						lv.Set(slog.LevelDebug)
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelWarn)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelError)).To(BeTrue())

						lv.Set(slog.LevelInfo)
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelWarn)).To(BeTrue())
						Ω(loggerz.Enabled(ctx, LevelError)).To(BeTrue())

						lv.Set(slog.LevelError)
						Ω(loggerz.Enabled(ctx, LevelDebug)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelInfo)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelWarn)).To(BeFalse())
						Ω(loggerz.Enabled(ctx, LevelError)).To(BeTrue())
					})
				})
			})
		})
		Convey("zap+", func() {
			core, obsLogs := observer.New(zap.DebugLevel)
			logger := zap.New(core).Sugar()

			getAndFormatZap := func() string {
				s := formatZap(obsLogs.All())
				fmt.Printf("\n%s\n", s)
				return s
			}
			Convey("zap", func() {
				assert := func() {
					s := getAndFormatZap()
					Ω(s).To(Equal("[INFO] Hello, World! name=Alice alias=A."))
				}

				Convey("zap.Infow", func() {
					logger.Infow("Hello, World!", "name", "Alice", "alias", "A.")

					assert()
				})

				loggerz := FromLoggerx(logger)
				Convey("loggerz.Infow", func() {
					loggerz.Infow("Hello, World!", "name", "Alice", "alias", "A.")

					assert()
				})
			})

			Convey("FromLoggerf", func() {
				logger := (*zapLoggerf)(logger)
				loggerz := FromLoggerf(logger)
				loggerz.Infow("Hello, World!", "name", "Alice", "alias", "A.")

				s := getAndFormatZap()
				fmt.Printf("\n%s\n", s)
				Ω(s).To(Equal(`[INFO] Hello, World! name="Alice" alias="A."`))
			})
		})
		Convey("log.Printf", func() {
			var b stringz.Builder
			logger := log.New(&b, "", 0)
			loggerz := FromLoggerP(logger)

			loggerz.Debugw("zero", "one", "1", "two", "2")
			loggerz.Infow("zero", "one", "1", "two", "2")
			loggerz.Infof("zero %v %v", "one", "two")
			loggerz.Warnf("zero %v %v", "one", "two")

			s := b.String()
			fmt.Printf("\n%s\n", s)
			expected := `
DEBUG: zero one="1" two="2"
 INFO: zero one="1" two="2"
 INFO: zero one two
 WARN: zero one two
`[1:]
			diffs := diffz.ByLine(s, expected)
			fmt.Println(diffs)
			Ω(diffs.IsDiff()).To(BeFalse())
		})
		Convey("keyValues", func() {
			var b stringz.Builder
			logger := log.New(&b, "", 0)
			loggerz := FromLoggerP(logger)

			loggerz.Debugw("zero", "one", "1", "two", "2")
			loggerz.Infow("extra", 0, "one", "two", 3, 4, "five")

			s := b.String()
			fmt.Printf("\n%s\n", s)
			expected := `
DEBUG: zero one="1" two="2"
 INFO: extra [0]=0 one="two" [1]=3 [2]=4 [3]="five"
`[1:]
			diffs := diffz.ByLine(s, expected)
			fmt.Println(diffs)
			Ω(diffs.IsDiff()).To(BeFalse())
		})
	})
}

type zapLoggerf zap.SugaredLogger

func (z *zapLoggerf) zap() *zap.SugaredLogger {
	return (*zap.SugaredLogger)(z)
}
func (z *zapLoggerf) Debugf(format string, args ...any) {
	z.zap().Debugf(format, args...)
}
func (z *zapLoggerf) Infof(format string, args ...any) {
	z.zap().Infof(format, args...)
}
func (z *zapLoggerf) Warnf(format string, args ...any) {
	z.zap().Warnf(format, args...)
}
func (z *zapLoggerf) Errorf(format string, args ...any) {
	z.zap().Errorf(format, args...)
}

func printLog(s string) {
	fmt.Printf("\n%s\n", strings.TrimSpace(s))
}

func formatZap(entries []observer.LoggedEntry) string {
	var b stringz.Builder
	for _, entry := range entries {
		b.Printf("[%v] %v", entry.Level.CapitalString(), entry.Message)
		for _, field := range entry.Context {
			b.Printf(" %v=%v", field.Key, field.String)
		}
		b.Println()
	}
	return strings.TrimSpace(b.String())
}
