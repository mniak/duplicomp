package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/cmd/sample/internal"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn := lo.Must(grpc.Dial(fmt.Sprintf(":%d", internal.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	))
	defer conn.Close()
	client := internal.NewPingerClient(conn)

	phrase := gofakeit.SentenceSimple()
	resp := lo.Must(client.SendPing(context.Background(), &internal.Ping{
		Message: &phrase,
	}))
	log.Printf("PONG %s", resp)
}
