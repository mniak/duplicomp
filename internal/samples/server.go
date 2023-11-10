package samples

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/mniak/ps121/internal/samples/grpc"
	g "google.golang.org/grpc"
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

	o.Logger.Print("Server Started. Waiting for calls.")
	server := g.NewServer()
	grpc.RegisterPingerServer(server, pinger)

	go func() {
		server.Serve(lis)
	}()
	return stoppable(func() {
		o.Logger.Print("Stopping gRPC server")
		server.GracefulStop()
		lis.Close()
		o.Logger.Print("Listener stopped")
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
