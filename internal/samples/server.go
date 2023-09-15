package samples

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/samples/grpc"
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

func defaultServerHandler(logger *log.Logger) _ServerHandler {
	return func(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
		meta, hasMeta := metadata.FromIncomingContext(ctx)
		logger.Printf("PING %s (hasMeta=%v, meta=%v)", *ping.Message, hasMeta, meta)

		return &grpc.Pong{
			OriginalMessage:    ping.Message,
			CapitalizedMessage: ptr(strings.ToUpper(*ping.Message)),
			RandomNumber:       ptr(gofakeit.Int32()),
		}, nil
	}
}

func (p _PingerServer) SendPing(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
	if p.options.ServerHandlerFactory == nil {
		return nil, errors.New("no handler defined")
	}
	handler := p.options.ServerHandlerFactory(p.options.Logger)
	return handler(ctx, ping)
}
