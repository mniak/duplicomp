package duplicomp

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Message struct {
	internalMessage proto.Message
}

type Stream interface {
	Send(m Message) error
	Receive() (*Message, error)
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

func (s *protoStream) Send(m Message) error {
	err := s.stream.SendMsg(m.internalMessage)
	return err
}

func (s *protoStream) Receive() (*Message, error) {
	msg := Message{
		internalMessage: new(emptypb.Empty),
	}
	err := s.stream.RecvMsg(msg.internalMessage)
	return &msg, err
}
