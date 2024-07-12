package logz

type xLogger struct {
	w  Loggerw
	f  Loggerf
	kv []any
}

func (l *xLogger) Debugw(msg string, keyVals ...any) {
	if l.w != nil {
		l.w.Debugw(msg, keyVals...)
	} else {
		l.f.Debugf("%v", formatWf(msg, "", keyVals))
	}
}
func (l *xLogger) Infow(msg string, keyVals ...any) {
	if l.w != nil {
		l.w.Infow(msg, keyVals...)
	} else {
		l.f.Infof("%v", formatWf(msg, "", keyVals))
	}
}
func (l *xLogger) Warnw(msg string, keyVals ...any) {
	if l.w != nil {
		l.w.Warnw(msg, keyVals...)
	} else {
		l.f.Warnf("%v", formatWf(msg, "", keyVals))
	}
}
func (l *xLogger) Errorw(msg string, keyVals ...any) {
	if l.w != nil {
		l.w.Errorw(msg, keyVals...)
	} else {
		l.f.Errorf("%v", formatWf(msg, "", keyVals))
	}
}
func (l *xLogger) Debugf(format string, args ...any) {
	if l.f != nil {
		l.f.Debugf(format, args...)
	} else {
		l.w.Debugw(formatf(format, "", args)())
	}
}
func (l *xLogger) Infof(format string, args ...any) {
	if l.f != nil {
		l.f.Infof(format, args...)
	} else {
		l.w.Infow(formatf(format, "", args)())
	}
}
func (l *xLogger) Warnf(format string, args ...any) {
	if l.f != nil {
		l.f.Warnf(format, args...)
	} else {
		l.w.Warnw(formatf(format, "", args)())
	}
}
func (l *xLogger) Errorf(format string, args ...any) {
	if l.f != nil {
		l.f.Errorf(format, args...)
	} else {
		l.w.Errorw(formatf(format, "", args)())
	}
}
func (l *xLogger) With(keyVals ...any) Logger {
	kv := append(l.kv[:], keyVals...)
	return &xLogger{w: l.w, f: l.f, kv: kv}
}
