package duplicomp

import "google.golang.org/protobuf/proto"

//go:generate mockgen -package=duplicomp -destination=mock_stream_output_test.go . OutputStream
type OutputStream interface {
	Send(m proto.Message) error
}
