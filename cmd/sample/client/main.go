package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGABRT)

	log.Print("Client started!")
	for {
		select {
		case <-sig:
		case <-time.After(5 * time.Second):
			phrase := gofakeit.SentenceSimple()
			resp, err := client.SendPing(context.Background(), &internal.Ping{
				Message: &phrase,
			})
			if err != nil {
				log.Printf("ERROR %s", err)
				continue
			}
			log.Printf("PONG %s", resp)
		}
	}
}
