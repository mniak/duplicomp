package ps121

import "google.golang.org/protobuf/proto"

//go:generate mockgen -package=ps121 -destination=mock_stream_output_test.go . OutputStream
type OutputStream interface {
	Send(m proto.Message) error
}
