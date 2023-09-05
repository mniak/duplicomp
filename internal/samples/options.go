package samples

import (
	"context"

	"github.com/mniak/duplicomp/internal/samples/internal"
)

type _ServerHandler func(ctx context.Context, ping *internal.Ping) (*internal.Pong, error)

type _Options struct {
	ServerHandler _ServerHandler
	Port          int
}
type _Option func(*_Options)

func WithHandler(h _ServerHandler) _Option {
	return func(o *_Options) {
		o.ServerHandler = h
	}
}

func WithPort(port int) _Option {
	return func(o *_Options) {
		o.Port = port
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
		ServerHandler: defaultServerHandler,
		Port:          9000,
	}
}
