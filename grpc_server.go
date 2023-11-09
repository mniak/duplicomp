package duplicomp

import (
	"context"
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
	server   *grpc.Server
}

func (self *GRPCServer) init() {
	if self.Logger == nil {
		self.Logger = noop.Logger()
	}
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
			self.Logger.Printf("Run failed with error: %s", self.runError)
		}
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

func (self *GRPCServer) GracefulStop() {
	if self.server != nil {
		self.server.GracefulStop()
		self.server = nil
	}
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

	serverStream := StreamsFromProtobuf(protoServer)
	return self.ConnectionHandler.HandleConnection(ctx, method, serverStream)
}
