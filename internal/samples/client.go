package samples

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/samples/grpc"
	"github.com/samber/lo"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func RunSendPing(phrase string, opts ..._Option) (*grpc.Pong, error) {
	o := defaultOptions().apply(opts...)
	conn := lo.Must(g.Dial(fmt.Sprintf(":%d", o.Port),
		g.WithTransportCredentials(insecure.NewCredentials()),
		g.WithUserAgent("sample-client/0.0.1"),
	))
	defer conn.Close()
	client := grpc.NewPingerClient(conn)
	meta := metadata.MD{
		"x-custom": []string{gofakeit.BuzzWord()},
	}
	ctx := metadata.NewOutgoingContext(context.Background(), meta)

	o.Logger.Printf("PING -- '%s'", phrase)
	pong, err := client.SendPing(ctx, &grpc.Ping{
		Message: &phrase,
	})
	panic("============> 267")
	if pong != nil {
		o.Logger.Printf("PONG -- '%s'", *pong.Reply)
	} else if err != nil {
		o.Logger.Printf("PONG ERROR -- %v", err)
	}
	return pong, err
}
