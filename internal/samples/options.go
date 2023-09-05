package samples

import (
	"context"

	"github.com/mniak/duplicomp/internal/samples/internal"
)

type _ServerHandler func(ctx context.Context, ping *internal.Ping) (*internal.Pong, error)

type _Options struct {
	ServerHandler _ServerHandler
}
type _Option func(*_Options)

func WithHandler(h _ServerHandler) _Option {
	return func(o *_Options) {
		o.ServerHandler = h
	}
}

func buildOptions(opts ..._Option) _Options {
	var options _Options
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
