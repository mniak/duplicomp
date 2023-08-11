package main

type InboundConfig struct {
	ListenAddress string
}

type OutboundConfig struct {
	TargetAddress  string
	IgnoreResponse bool
}

type Proxy struct {
	InboundConfig   InboundConfig
	OutboundConfigs []OutboundConfig
}
