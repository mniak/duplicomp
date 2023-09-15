package duplicomp

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *protoStream) Receive() (proto.Message, error) {
	msg := new(emptypb.Empty)
	err := s.stream.RecvMsg(msg)
	return msg, err
}
