package samples

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/samples/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RunServer(port int, opts ..._Option) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	defer lis.Close()

	pinger := _PingerServer{
		options: buildOptions(opts...),
	}

	log.Println("Server Started. Waiting for calls.")
	server := grpc.NewServer()
	internal.RegisterPingerServer(server, pinger)

	return server.Serve(lis)
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
	log.Printf("PING %s (hasMeta=%v, meta=%v)", *ping.Message, hasMeta, meta)

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

func RunMockServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	defer lis.Close()

	pinger := _PingerServer{}

	log.Println("Server Started. Waiting for calls.")
	server := grpc.NewServer()
	internal.RegisterPingerServer(server, pinger)

	return server.Serve(lis)
}
