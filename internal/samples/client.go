package samples

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/samples/internal"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func RunSendPing(opts ..._Option) error {
	o := defaultOptions().apply(opts...)
	conn := lo.Must(grpc.Dial(fmt.Sprintf(":%d", o.Port),
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
	o.Logger.Printf("PING %s", phrase)
	resp, err := client.SendPing(ctx, &internal.Ping{
		Message: &phrase,
	})
	if err != nil {
		o.Logger.Print("ERROR %s", err)
		return err
	}
	o.Logger.Printf("PONG %s", resp)
	return nil
}
