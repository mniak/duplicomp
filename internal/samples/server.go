package samples

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/mniak/duplicomp/internal/samples/grpc"
	"github.com/mniak/duplicomp/log2"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RunServer(opts ..._Option) (stop Stoppable, err error) {
	o := *defaultOptions().apply(opts...)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", o.Port))
	if err != nil {
		return nil, err
	}
	defer func() {
		if stop == nil {
			lis.Close()
		}
	}()

	pinger := _PingerServer{
		options: o,
	}

	o.Logger.Println("Server Started. Waiting for calls.")
	server := g.NewServer()
	grpc.RegisterPingerServer(server, pinger)

	go func() {
		server.Serve(lis)
	}()
	return stoppable(func() {
		o.Logger.Println("Stopping gRPC server")
		server.GracefulStop()
		lis.Close()
		o.Logger.Println("Listener stopped")
	}), nil
}

type _PingerServer struct {
	grpc.PingerServer
	options _Options
}

func ptr[T any](t T) *T {
	return &t
}

func (p _PingerServer) SendPing(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
	if p.options.ServerHandler == nil {
		return nil, errors.New("no handler defined")
	}
	return p.options.ServerHandler.Handle(ctx, ping)
}

type _ServerHandler func(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error)

func (h _ServerHandler) Handle(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
	return h(ctx, ping)
}

type ServerHandler interface {
	Handle(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error)
}

type defaultServerHandler struct {
	logger log2.Logger
}

func (h defaultServerHandler) Handle(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
	meta, hasMeta := metadata.FromIncomingContext(ctx)
	h.logger.Printf("PING %s (hasMeta=%v, meta=%v)", *ping.Message, hasMeta, meta)

	return &grpc.Pong{
		Reply:    ping.Message,
		ServedBy: ptr("Default Server Handler"),
	}, nil
}
