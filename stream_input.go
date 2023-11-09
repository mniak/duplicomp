package duplicomp

import "google.golang.org/protobuf/proto"

//go:generate mockgen -package=duplicomp -destination=mock_stream_input_test.go . InputStream
type InputStream interface {
	Receive() (proto.Message, error)
}
