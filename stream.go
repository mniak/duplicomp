package ps121

import (
	"context"

	"github.com/mniak/ps121/empty"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

//go:generate mockgen -package=ps121 -destination=mock_stream_test.go . Stream
type Stream interface {
	InputStream
	OutputStream
}

type inOutStream struct {
	InputStream
	OutputStream
	method string
}

func InOutStream(in InputStream, out OutputStream) Stream {
	return &inOutStream{
		InputStream:  in,
		OutputStream: out,
	}
}

type iProtoStream interface {
	Context() context.Context
	SendMsg(m interface{}) error
	RecvMsg(m interface{}) error
}

func StreamsFromProtobuf(s iProtoStream) Stream {
	pstr := protoStream{
		stream: s,
	}
	return InOutStream(&pstr, &pstr)
}

type protoStream struct {
	stream iProtoStream
	method string
}

func (s *protoStream) Send(m proto.Message) error {
	err := s.stream.SendMsg(m)
	return err
}

type MyMessageType struct{}

func buildDescriptor() protoreflect.MessageDescriptor {
	return protoimpl.X.MessageDescriptorOf(&empty.Empty{})
}

func (s *protoStream) Receive() (proto.Message, error) {
	// msg := dynamicpb.NewMessage(buildDescriptor())
	msg := new(empty.Empty)
	err := s.stream.RecvMsg(msg)
	return msg, err
}
