package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/cmd/sample/internal"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	var port int
	flag.IntVar(&port, "port", internal.Port, "TCP port to connect")
	flag.Parse()

	conn := lo.Must(grpc.Dial(fmt.Sprintf(":%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent("sample-client/0.0.1"),
	))
	defer conn.Close()
	client := internal.NewPingerClient(conn)

	meta := metadata.MD{
		"x-custom": []string{gofakeit.BuzzWord()},
	}
	ctx := metadata.NewOutgoingContext(context.Background(), meta)

	phrase := gofakeit.SentenceSimple()
	resp, err := client.SendPing(ctx, &internal.Ping{
		Message: &phrase,
	})
	if err != nil {
		log.Printf("ERROR %s", err)
		return
	}
	log.Printf("PONG %s", resp)
}
