package log2

import "testing"

type testingLogger struct {
	t *testing.T
}

func (l *testingLogger) Printf(format string, args ...interface{}) {
	l.t.Helper()
	l.t.Logf(format, args...)
}

func (l *testingLogger) Println(args ...interface{}) {
	l.t.Helper()
	l.t.Log(args...)
}

func FromT(t *testing.T) Logger {
	t.Helper()
	return &testingLogger{t}
}
