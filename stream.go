package duplicomp

import (
	"github.com/mniak/duplicomp/empty"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type Stream interface {
	InputStream
	OutputStream
}

type inOutStream struct {
	InputStream
	OutputStream
}

func InOutStream(in InputStream, out OutputStream) Stream {
	return &inOutStream{
		InputStream:  in,
		OutputStream: out,
	}
}

type iProtoStream interface {
	SendMsg(m interface{}) error
	RecvMsg(m interface{}) error
}

func StreamFromProtobuf(s iProtoStream) (InputStream, OutputStream) {
	str := protoStream{
		stream: s,
	}
	return &str, &str
}

type protoStream struct {
	stream iProtoStream
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
