package samples

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/samples/internal"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func RunSendPing(opts ..._Option) error {
	options := *defaultOptions().apply(opts...)
	conn := lo.Must(grpc.Dial(fmt.Sprintf(":%d", options.Port),
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
	log.Printf("PING %s", phrase)
	resp, err := client.SendPing(ctx, &internal.Ping{
		Message: &phrase,
	})
	if err != nil {
		log.Printf("ERROR %s", err)
		return err
	}
	log.Printf("PONG %s", resp)
	return nil
}
