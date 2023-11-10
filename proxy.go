package ps121

import "context"

type InboundConfig struct {
	ListenAddress string
}

type OutboundConfig struct {
	TargetAddress  string
	IgnoreResponse bool
}

type ProxyConfig struct {
	InboundConfig   InboundConfig
	OutboundConfigs []OutboundConfig
}

type Proxy interface {
	Start() error
	Stop()
	Wait() error
}

type TargetConnection interface {
	Stream(ctx context.Context, method string) (Stream, error)
	Close()
}
