package ps121

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func copyHeadersFromIncomingToOutcoming(in, out context.Context) context.Context {
	meta, hasMeta := metadata.FromIncomingContext(in)
	if hasMeta {
		for k, vals := range meta {
			for _, v := range vals {
				in = metadata.AppendToOutgoingContext(out, k, v)
			}
		}
	}
	return in
}
