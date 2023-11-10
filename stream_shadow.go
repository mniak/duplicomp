package ps121

import (
	"io"
	"sync"
	"time"

	"github.com/mniak/ps121/internal/noop"
	"github.com/mniak/ps121/log2"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type MessageBytesExtractor interface {
	ExtractBytes(proto.Message) []byte
}

type SimpleMessageBytesExtractor struct{}

func (SimpleMessageBytesExtractor) ExtractBytes(m proto.Message) []byte {
	return m.ProtoReflect().GetUnknown()
}

type StreamWithShadow struct {
	Method                string
	Primary               Stream
	Shadow                Stream
	Logger                log2.Logger
	Comparator            Comparator
	MessageBytesExtractor MessageBytesExtractor
	BufferSize            int

	onceInit                sync.Once
	shadowInputChan         OverflowChannel[ReceivedMessage]
	onceStartShadowReceiver sync.Once
}

func (self *StreamWithShadow) init() {
	self.onceInit.Do(func() {
		if self.Logger == nil {
			self.Logger = noop.Logger()
		}
		if self.Comparator == nil {
			self.Comparator = noop.Comparator()
		}
		self.shadowInputChan = NewOverflowableChannel[ReceivedMessage](self.BufferSize)
	})
}

func (self *StreamWithShadow) Send(m proto.Message) error {
	self.init()

	err := self.Primary.Send(m)
	if err != nil {
		return err
	}

	go func() {
		err := self.Shadow.Send(m)
		if err != nil {
			self.Logger.Printf("failed sending to shadow: %v", err)
		}
	}()

	return nil
}

type ReceivedMessage struct {
	Message proto.Message
	Error   error
	Time    time.Time
}

func (self *StreamWithShadow) shadowReceiveLoop() {
	for primaryMsg := range self.shadowInputChan.Receiver() {

		shadowMsg, shadowErr := self.Shadow.Receive()

		var msgBytes []byte
		if primaryMsg.Message != nil {
			msgBytes = primaryMsg.Message.ProtoReflect().GetUnknown()
		}
		var shadowMsgBytes []byte
		if shadowMsg != nil {
			shadowMsgBytes = shadowMsg.ProtoReflect().GetUnknown()
		}

		compareError := self.Comparator.Compare(
			self.Method,
			msgBytes, primaryMsg.Error,
			shadowMsgBytes, shadowErr,
		)

		if compareError != nil {
			self.Logger.Print(errors.WithMessage(compareError, "comparison failed"))
		}
	}
}

func (self *StreamWithShadow) Receive() (proto.Message, error) {
	self.init()

	msg, err := self.Primary.Receive()

	self.onceStartShadowReceiver.Do(func() {
		go self.shadowReceiveLoop()
	})

	self.shadowInputChan.Send(ReceivedMessage{
		Message: msg,
		Error:   err,
		Time:    time.Now(),
	})

	if err == io.EOF {
		self.shadowInputChan.Close()
	}

	return msg, err
}

const OverflowableChannelDefaultBufferSize = 2
