package duplicomp

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

// func (p *GRPCServer) LoggingTo(w io.Writer) *GRPCServer {
// 	return p.LoggingToWithPrefix(w, "[Proxy] ")
// }

// func (p *GRPCServer) LoggingToWithPrefix(w io.Writer, prefix string) *GRPCServer {
// 	p.logger = log.New(w, prefix, 0)
// 	return p
// }

// func (p *GRPCServer) WithLogger(logger log2.Logger) *GRPCServer {
// 	p.logger = logger
// 	return p
// }

type TargetConnection interface {
	Stream() Stream
	Close()
}

func connectTo(outConfig OutboundConfig) TargetConnection {
	// for _, outConfig := range outConfigs {
	// 	clientCredentials := insecure.NewCredentials()
	// 	conn, err := grpc.Dial(outConfig.TargetAddress,
	// 		grpc.WithTransportCredentials(clientCredentials),
	// 		grpc.WithUserAgent("duplicomp-gateway/0.0.1"),
	// 	)
	// 	coinn
	// 	if err != nil {
	// 		break
	// 	}
	// 	p.connections = append(p.connections, conn)
	// }
	// if err != nil {
	// 	p.Stop()
	// 	return err
	// }
	return nil
}
