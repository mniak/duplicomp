package samples

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/samples/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var defaultServerLogger = log.New(os.Stdout, "[Server] ", 0)

func RunServer(opts ..._Option) (stop Stoppable, err error) {
	logger := defaultServerLogger

	options := *defaultOptions().apply(opts...)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", options.Port))
	if err != nil {
		return nil, err
	}
	defer func() {
		if stop == nil {
			lis.Close()
		}
	}()

	pinger := _PingerServer{
		options: options,
	}

	logger.Println("Server Started. Waiting for calls.")
	server := grpc.NewServer()
	internal.RegisterPingerServer(server, pinger)

	go func() {
		server.Serve(lis)
	}()
	return stoppable(func() {
		logger.Println("Stopping gRPC server")
		server.GracefulStop()
		lis.Close()
		logger.Println("Listener stopped")
	}), nil
}

type _PingerServer struct {
	internal.PingerServer
	options _Options
}

func ptr[T any](t T) *T {
	return &t
}

func defaultServerHandler(ctx context.Context, ping *internal.Ping) (*internal.Pong, error) {
	meta, hasMeta := metadata.FromIncomingContext(ctx)
	defaultServerLogger.Printf("PING %s (hasMeta=%v, meta=%v)", *ping.Message, hasMeta, meta)

	return &internal.Pong{
		OriginalMessage:    ping.Message,
		CapitalizedMessage: ptr(strings.ToUpper(*ping.Message)),
		RandomNumber:       ptr(gofakeit.Int32()),
	}, nil
}

func (p _PingerServer) SendPing(ctx context.Context, ping *internal.Ping) (*internal.Pong, error) {
	if p.options.ServerHandler == nil {
		return nil, errors.New("no handler defined")
	}
	return p.options.ServerHandler(ctx, ping)
}
