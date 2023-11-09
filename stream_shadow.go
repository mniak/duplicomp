package duplicomp

import (
	"io"
	"sync"
	"time"

	"github.com/mniak/duplicomp/internal/noop"
	"github.com/mniak/duplicomp/log2"
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
	Primary               Stream
	Shadow                Stream
	Logger                log2.Logger
	Comparator            Comparator
	MessageBytesExtractor MessageBytesExtractor
	BufferSize            int

	onceInit                sync.Once
	shadowInputChan         OverflowableChannel[ReceivedMessage]
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

type OverflowableChannel[T any] struct {
	dataChan     chan T
	closedSignal chan struct{}
	onceClose    sync.Once
}

const OverflowableChannelDefaultBufferSize = 2

func NewOverflowableChannel[T any](bufferSize int) OverflowableChannel[T] {
	if bufferSize < 1 {
		bufferSize = OverflowableChannelDefaultBufferSize
	}
	return OverflowableChannel[T]{
		dataChan:     make(chan T, bufferSize),
		closedSignal: make(chan struct{}),
	}
}

func (self OverflowableChannel[T]) Send(val T) {
	select {
	case <-self.closedSignal:

	default:
		select {
		case self.dataChan <- val:
		default:
			self.Close()
		}

	}
}

func (self OverflowableChannel[T]) Close() {
	self.onceClose.Do(func() {
		close(self.closedSignal)
		close(self.dataChan)
	})
}

func (self OverflowableChannel[T]) Receiver() <-chan T {
	return self.dataChan
}
