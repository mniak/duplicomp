package duplicomp

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/mniak/duplicomp/internal/noop"
	"github.com/mniak/duplicomp/log2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

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

type GRPCProxy struct {
	config ProxyConfig
	logger log2.Logger

	listener    net.Listener
	connections []*grpc.ClientConn
	runError    error
	stopped     chan struct{}
	server      *grpc.Server
}

func NewGRPCProxy(config ProxyConfig) *GRPCProxy {
	return &GRPCProxy{
		config:  config,
		stopped: make(chan struct{}),
		logger:  noop.Logger(),
	}
}

func (p *GRPCProxy) LoggingTo(w io.Writer) *GRPCProxy {
	return p.LoggingToWithPrefix(w, "[Proxy] ")
}

func (p *GRPCProxy) LoggingToWithPrefix(w io.Writer, prefix string) *GRPCProxy {
	p.logger = log.New(w, prefix, 0)
	return p
}

func (p *GRPCProxy) WithLogger(logger log2.Logger) *GRPCProxy {
	p.logger = logger
	return p
}

func (p *GRPCProxy) run() error {
	p.server = grpc.NewServer(grpc.UnknownServiceHandler(p.connectionHandler))
	p.server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "DummyService",
		HandlerType: (*any)(nil),
	}, nil)

	err := p.server.Serve(p.listener)
	if err != nil {
		return err
	}
	return nil
}

func (p *GRPCProxy) connectionHandler(_ any, protoServer grpc.ServerStream) error {
	ctx := protoServer.Context()
	ctx, stop := context.WithCancel(ctx)
	defer stop()

	method, hasName := grpc.MethodFromServerStream(protoServer)
	if !hasName {
		return status.Errorf(codes.NotFound, "Method name could not be determined")
	}
	p.logger.Printf("Handling method %s", method)
	defer p.logger.Printf("Done handling method %s", method)

	ctx = copyHeadersFromIncomingToOutcoming(ctx, ctx)

	// server := NewProtoStream(protoServer)
	// serverObservable := ObservableFromStream(ctx, server)

	// if p.UseShadow {
	// 	go func() {
	// 		shadow := Forwarder{
	// 			LogName:           "Shadow",
	// 			Method:            method,
	// 			Server:            server,
	// 			ServerObservable:  serverObservable,
	// 			InboundConnection: p.ShadowClientConnection,
	// 			DiscardResponses:  true,
	// 		}
	// 		err := shadow.Run(ctx)
	// 		if err != nil {
	// 			log.Printf("[Shadow] Failure: %s", err.Error())
	// 		}
	// 	}()
	// }

	// primary := Forwarder{
	// 	Method:            method,
	// 	Server:            server,
	// 	ServerObservable:  serverObservable,
	// 	InboundConnection: p.PrimaryClientConnection,
	// }
	// err := primary.Run(ctx)
	// if err != nil {
	// 	return status.Errorf(codes.Internal, err.Error())
	// }
	return nil
}

func (p *GRPCProxy) runAsync() {
	p.runError = nil
	go func() {
		p.runError = p.run()
		if p.runError != nil {
			log.Printf("[Proxy] run failed with error: %s", p.runError)
		}
		close(p.stopped)
	}()
}

func (p *GRPCProxy) Start() error {
	var err error
	p.listener, err = net.Listen("tcp", p.config.InboundConfig.ListenAddress)
	if err != nil {
		return err
	}
	for _, outConfig := range p.config.OutboundConfigs {
		clientCredentials := insecure.NewCredentials()
		conn, err := grpc.Dial(outConfig.TargetAddress,
			grpc.WithTransportCredentials(clientCredentials),
			grpc.WithUserAgent("duplicomp-gateway/0.0.1"),
		)
		if err != nil {
			break
		}
		p.connections = append(p.connections, conn)
	}
	if err != nil {
		p.Stop()
		return err
	}
	p.runAsync()
	return nil
}

func (p *GRPCProxy) Stop() {
	if p.listener != nil {
		p.listener.Close()
		p.listener = nil
	}

	for _, conn := range p.connections {
		conn.Close()
	}
}

func (p *GRPCProxy) Wait() error {
	if p.stopped != nil {
		<-p.stopped
	}
	return p.runError
}
