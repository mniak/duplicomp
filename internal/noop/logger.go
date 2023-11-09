package noop

import (
	"io"
	"log"

	"github.com/mniak/duplicomp/log2"
)

func Logger() log2.Logger {
	return log.New(io.Discard, "", log.LstdFlags)
}
