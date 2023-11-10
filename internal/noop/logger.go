package noop

import (
	"io"
	"log"

	"github.com/mniak/ps121/log2"
)

func Logger() log2.Logger {
	return log.New(io.Discard, "", log.LstdFlags)
}
