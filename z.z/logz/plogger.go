package logz

type pLogger struct {
	l  LoggerP
	kv []any
}

func (l *pLogger) Debugw(msg string, keyValues ...any) {
	l.l.Printf("%v", formatWf(msg, strDebug, keyValues))
}
func (l *pLogger) Infow(msg string, keyValues ...any) {
	l.l.Printf("%v", formatWf(msg, strInfo, keyValues))
}
func (l *pLogger) Warnw(msg string, keyValues ...any) {
	l.l.Printf("%v", formatWf(msg, strWarn, keyValues))
}
func (l *pLogger) Errorw(msg string, keyValues ...any) {
	l.l.Printf("%v", formatWf(msg, strError, keyValues))
}
func (l *pLogger) Debugf(format string, args ...any) {
	l.l.Printf("%v", formatf(format, strDebug, args))
}
func (l *pLogger) Infof(format string, args ...any) {
	l.l.Printf("%v", formatf(format, strInfo, args))
}
func (l *pLogger) Warnf(format string, args ...any) {
	l.l.Printf("%v", formatf(format, strWarn, args))
}
func (l *pLogger) Errorf(format string, args ...any) {
	l.l.Printf("%v", formatf(format, strError, args))
}
func (l *pLogger) With(keyValues ...any) Logger {
	cloned := *l
	l.kv = append(l.kv[:], keyValues...)
	return &cloned
}
