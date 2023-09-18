package duplicomp

import (
	"context"
	"log"
	"net"

	"github.com/mniak/duplicomp/internal/noop"
	"github.com/mniak/duplicomp/log2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConnectionHandler interface {
	HandleConnection(ctx context.Context, method string, serverStream Stream) error
}

type GRPCServer struct {
	ConnectionHandler ConnectionHandler
	Logger            log2.Logger

	listener net.Listener
	runError error
	stopped  chan struct{}
	server   *grpc.Server
}

func (p *GRPCServer) init() {
	if p.Logger == nil {
		p.Logger = noop.Logger()
	}
	p.stopped = make(chan struct{})
}

func (p *GRPCServer) Start(addr string) error {
	var err error
	p.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return p.StartWithListener(p.listener)
}

func (p *GRPCServer) StartWithListener(lis net.Listener) error {
	p.init()
	p.listener = lis
	p.runAsync()
	return nil
}

func (p *GRPCServer) runAsync() {
	p.runError = nil
	go func() {
		p.runError = p.run()
		if p.runError != nil {
			log.Printf("[Proxy] run failed with error: %s", p.runError)
		}
		close(p.stopped)
	}()
}

func (p *GRPCServer) run() error {
	p.server = grpc.NewServer(grpc.UnknownServiceHandler(p.onConnectionAccepted))
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

func (p *GRPCServer) Stop() {
	if p.listener != nil {
		p.listener.Close()
		p.listener = nil
	}
}

func (p *GRPCServer) Wait() error {
	if p.stopped != nil {
		<-p.stopped
	}
	return p.runError
}

func (p *GRPCServer) onConnectionAccepted(_ any, protoServer grpc.ServerStream) error {
	ctx := protoServer.Context()

	method, hasName := grpc.MethodFromServerStream(protoServer)
	if !hasName {
		return status.Errorf(codes.NotFound, "Method name could not be determined")
	}
	p.Logger.Printf("Handling method %s", method)
	defer p.Logger.Printf("Done handling method %s", method)

	ctx = copyHeadersFromIncomingToOutcoming(ctx, ctx)

	serverStream := InOutStream(StreamFromProtobuf(protoServer))
	// serverObservable := ObservableFromStream(ctx, server)

	return p.ConnectionHandler.HandleConnection(ctx, method, serverStream)
}
