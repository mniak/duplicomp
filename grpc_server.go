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

func (self *GRPCServer) init() {
	if self.Logger == nil {
		self.Logger = noop.Logger()
	}
	self.stopped = make(chan struct{})
}

func (p *GRPCServer) Start(addr string) error {
	var err error
	p.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return p.StartWithListener(p.listener)
}

func (self *GRPCServer) StartWithListener(lis net.Listener) error {
	self.init()
	self.listener = lis
	self.runAsync()
	return nil
}

func (self *GRPCServer) runAsync() {
	self.runError = nil
	go func() {
		self.runError = self.run()
		if self.runError != nil {
			log.Printf("[Proxy] run failed with error: %s", self.runError)
		}
		close(self.stopped)
	}()
}

func (self *GRPCServer) run() error {
	self.server = grpc.NewServer(grpc.UnknownServiceHandler(self.onConnectionAccepted))
	self.server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "DummyService",
		HandlerType: (*any)(nil),
	}, nil)

	err := self.server.Serve(self.listener)
	if err != nil {
		return err
	}
	return nil
}

func (self *GRPCServer) Stop() {
	if self.listener != nil {
		self.listener.Close()
		self.listener = nil
	}
}

func (self *GRPCServer) Wait() error {
	if self.stopped != nil {
		<-self.stopped
	}
	return self.runError
}

func (self *GRPCServer) onConnectionAccepted(_ any, protoServer grpc.ServerStream) error {
	ctx := protoServer.Context()

	method, hasName := grpc.MethodFromServerStream(protoServer)
	if !hasName {
		return status.Errorf(codes.NotFound, "Method name could not be determined")
	}
	self.Logger.Printf("Handling method %s", method)
	defer self.Logger.Printf("Done handling method %s", method)

	ctx = copyHeadersFromIncomingToOutcoming(ctx, ctx)

	serverStream := InOutStream(StreamFromProtobuf(protoServer))
	// serverObservable := ObservableFromStream(ctx, server)

	return self.ConnectionHandler.HandleConnection(ctx, method, serverStream)
}
