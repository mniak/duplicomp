package samples

import (
	"fmt"
	"log"
	"os"

	"github.com/mniak/duplicomp/internal/noop"
	"github.com/mniak/duplicomp/log2"
)

type _Options struct {
	Port          int
	Logger        log2.Logger
	ServerHandler ServerHandler
}
type _Option func(*_Options)

func WithHandler(h ServerHandler) _Option {
	return func(o *_Options) {
		o.ServerHandler = h
	}
}

func WithPort(port int) _Option {
	return func(o *_Options) {
		o.Port = port
	}
}

func WithName(name string) _Option {
	return func(o *_Options) {
		o.Logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", name), log.LstdFlags)
	}
}

func WithLogger(logger log2.Logger) _Option {
	return func(o *_Options) {
		o.Logger = logger
		if sh, is := o.ServerHandler.(defaultServerHandler); is {
			sh.logger = logger
		}
	}
}

func (o *_Options) apply(opts ..._Option) *_Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func buildOptions(opts ..._Option) _Options {
	var options _Options
	options.apply(opts...)
	return options
}

func defaultOptions() *_Options {
	return &_Options{
		Port:          9000,
		Logger:        noop.Logger(),
		ServerHandler: defaultServerHandler{logger: noop.Logger()},
	}
}
