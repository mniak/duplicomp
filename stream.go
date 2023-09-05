package duplicomp

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Stream interface {
	Send(m proto.Message) error
	Receive() (proto.Message, error)
}

type iProtoStream interface {
	SendMsg(m interface{}) error
	RecvMsg(m interface{}) error
}

func NewProtoStream(s iProtoStream) Stream {
	return &protoStream{
		stream: s,
	}
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
	err := s.stream.RecvMsg(&msg)
	return msg, err
}
