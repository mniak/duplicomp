package main

import (
	"google.golang.org/protobuf/types/known/emptypb"
)

type Message struct {
	internalMessage emptypb.Empty
}

type Stream interface {
	SendMsg(m Message) error
	RecvMsg() (Message, error)
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

func (s *protoStream) SendMsg(m Message) error {
	err := s.stream.SendMsg(&m.internalMessage)
	return err
}

func (s *protoStream) RecvMsg() (Message, error) {
	var msg Message
	err := s.stream.RecvMsg(&msg.internalMessage)
	return msg, err
}
