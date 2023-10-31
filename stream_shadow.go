package duplicomp

import (
	"sync"

	"github.com/mniak/duplicomp/internal/noop"
	"github.com/mniak/duplicomp/log2"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type StreamWithShadow struct {
	Primary    Stream
	Shadow     Stream
	Logger     log2.Logger
	Comparator Comparator

	once sync.Once
}

func (self *StreamWithShadow) init() {
	self.once.Do(func() {
		if self.Logger == nil {
			self.Logger = noop.Logger()
		}
		if self.Comparator == nil {
			self.Comparator = noop.Comparator()
		}
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

func (self *StreamWithShadow) Receive() (proto.Message, error) {
	msg, err := self.Primary.Receive()

	go func() {
		shadowMsg, shadowErr := self.Shadow.Receive()

		msgBytes := msg.ProtoReflect().GetUnknown()
		shadowMsgBytes := shadowMsg.ProtoReflect().GetUnknown()

		compareError := self.Comparator.Compare(
			msgBytes, err,
			shadowMsgBytes, shadowErr,
		)

		if compareError != nil {
			self.Logger.Print(errors.WithMessage(compareError, "comparison failed"))
		}
	}()

	return msg, err
}
