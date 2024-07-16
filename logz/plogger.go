package logz

type pLogger struct {
	l  LoggerP
	kv []zKV
	fn EnablerFunc
}

func (l *pLogger) Debugw(msg string, keyValues ...any) {
	l.l.Printf("%v", formatWf(strDebug, msg, l.kv, keyValues))
}
func (l *pLogger) Infow(msg string, keyValues ...any) {
	l.l.Printf("%v", formatWf(strInfo, msg, l.kv, keyValues))
}
func (l *pLogger) Warnw(msg string, keyValues ...any) {
	l.l.Printf("%v", formatWf(strWarn, msg, l.kv, keyValues))
}
func (l *pLogger) Errorw(msg string, keyValues ...any) {
	l.l.Printf("%v", formatWf(strError, msg, l.kv, keyValues))
}
func (l *pLogger) Debugf(format string, args ...any) {
	l.l.Printf("%v", formatf(strDebug, format, args, l.kv))
}
func (l *pLogger) Infof(format string, args ...any) {
	l.l.Printf("%v", formatf(strInfo, format, args, l.kv))
}
func (l *pLogger) Warnf(format string, args ...any) {
	l.l.Printf("%v", formatf(strWarn, format, args, l.kv))
}
func (l *pLogger) Errorf(format string, args ...any) {
	l.l.Printf("%v", formatf(strError, format, args, l.kv))
}
func (l *pLogger) Enabled(level Level) bool {
	return l.fn(level)
}
func (l *pLogger) With(keyValues ...any) Logger {
	cloned := *l
	l.kv = appendKV(l.kv, keyValues)
	return &cloned
}
