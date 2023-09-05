package samples

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/samples/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RunServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	defer lis.Close()

	pinger := Pinger{}

	log.Println("Server Started. Waiting for calls.")
	server := grpc.NewServer()
	internal.RegisterPingerServer(server, pinger)
	return server.Serve(lis)
}

type Pinger struct {
	internal.PingerServer
}

func ptr[T any](t T) *T {
	return &t
}

func (p Pinger) SendPing(ctx context.Context, ping *internal.Ping) (*internal.Pong, error) {
	meta, hasMeta := metadata.FromIncomingContext(ctx)
	log.Printf("PING %s (hasMeta=%v, meta=%v)", *ping.Message, hasMeta, meta)

	return &internal.Pong{
		OriginalMessage:    ping.Message,
		CapitalizedMessage: ptr(strings.ToUpper(*ping.Message)),
		RandomNumber:       ptr(gofakeit.Int32()),
	}, nil
}
