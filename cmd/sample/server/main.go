package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/cmd/sample/internal"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

func main() {
	lis := lo.Must(net.Listen("tcp", fmt.Sprintf("localhost:%d", internal.Port)))
	defer lis.Close()

	pinger := Pinger{}

	server := grpc.NewServer()
	internal.RegisterPingerServer(server, pinger)
	server.Serve(lis)
}

type Pinger struct {
	internal.PingerServer
}

func ptr[T any](t T) *T {
	return &t
}

func (p Pinger) SendPing(ctx context.Context, ping *internal.Ping) (*internal.Pong, error) {
	log.Printf("PING %s", ping.Message)

	return &internal.Pong{
		OriginalMessage:    ping.Message,
		CapitalizedMessage: ptr(strings.ToUpper(*ping.Message)),
		RandomNumber:       ptr(gofakeit.Int32()),
	}, nil
}
