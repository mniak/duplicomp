package log2

type Logger interface {
	Printf(string, ...any)
	Print(...any)
}
