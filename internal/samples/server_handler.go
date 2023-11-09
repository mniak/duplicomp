package samples

import (
	"context"
	"fmt"

	"github.com/mniak/duplicomp/internal/samples/grpc"
	"github.com/mniak/duplicomp/log2"
	"google.golang.org/grpc/metadata"
)

//go:generate mockgen -package=samples -destination=mock_server_handler.go . ServerHandler
type ServerHandler interface {
	Handle(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error)
}

type defaultServerHandler struct {
	logger log2.Logger
}

func (h defaultServerHandler) Handle(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
	meta, hasMeta := metadata.FromIncomingContext(ctx)
	h.logger.Printf("PING %s (hasMeta=%v, meta=%v)", *ping.Message, hasMeta, meta)

	return &grpc.Pong{
		Reply:    ping.Message,
		ServedBy: ptr(fmt.Sprintf("Default Server Handler v%s", ApplicationVersion)),
	}, nil
}
