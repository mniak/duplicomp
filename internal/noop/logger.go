package noop

import (
	"io"
	"log"
)

func Logger() *log.Logger {
	return log.New(io.Discard, "", 0)
}
