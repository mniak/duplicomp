package duplicomp

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type StreamWithShadow struct {
	Primary Stream
	Shadow  Stream
	Logger  ShadowLogger
}

func (self *StreamWithShadow) Send(m proto.Message) error {
	err := self.Primary.Send(m)
	if err != nil {
		return err
	}

	go func() {
		err := self.Shadow.Send(m)
		if err != nil {
			self.Logger.LogSendFailure(err)
		}
	}()

	return nil
}

func (self *StreamWithShadow) Receive() (proto.Message, error) {
	msg, err := self.Primary.Receive()
	if err != nil {
		return msg, err
	}

	go func() {
		shadowMsg, shadowErr := self.Shadow.Receive()
		_, _ = shadowMsg, shadowErr
		fmt.Println("Shadow logger", self.Logger)
		self.Logger.LogCompareReceive(msg, shadowMsg, shadowErr)
	}()

	return msg, nil
}
