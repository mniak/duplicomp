package log2

import (
	"fmt"
	"io"
)

type writerLogger struct {
	w io.Writer
}

func (l *writerLogger) Printf(format string, args ...interface{}) {
	fmt.Fprintf(l.w, format+"\n", args...)
}

func (l *writerLogger) Print(args ...interface{}) {
	fmt.Fprintln(l.w, args...)
}

func FromWriter(t io.Writer) Logger {
	return &writerLogger{t}
}
