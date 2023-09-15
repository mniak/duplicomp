package log2

type Logger interface {
	Printf(string, ...any)
	Println(...any)
}
