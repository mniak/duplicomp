package log2

import "fmt"

type subLogger struct {
	logger Logger
	prefix string
}

func (l *subLogger) Printf(format string, args ...interface{}) {
	l.logger.Println(l.prefix + fmt.Sprintf(format, args...))
}

func (l *subLogger) Println(args ...interface{}) {
	l.logger.Println(l.prefix + fmt.Sprint(args...))
}

func Sub(logger Logger, prefix string) Logger {
	return &subLogger{
		logger: logger,
		prefix: prefix,
	}
}
