package duplicomp

import (
	"google.golang.org/protobuf/proto"
)

//go:generate mockgen -package=duplicomp -destination=mock_shadowlogger_test.go . ShadowLogger
type ShadowLogger interface {
	LogSendFailure(error)
	LogCompareReceive(primaryMsg, shadowMsg proto.Message, shadowErr error)
}

type StreamWithShadow struct {
	Primary Stream
	Shadow  Stream
	Logger  ShadowLogger
}

func (fs *StreamWithShadow) Send(m proto.Message) error {
	err := fs.Primary.Send(m)
	if err != nil {
		return err
	}

	go func() {
		err := fs.Shadow.Send(m)
		if err != nil {
			fs.Logger.LogSendFailure(err)
		}
	}()

	return nil
}

func (fs *StreamWithShadow) Receive() (proto.Message, error) {
	msg, err := fs.Primary.Receive()
	if err != nil {
		return msg, err
	}

	go func() {
		shadowMsg, shadowErr := fs.Shadow.Receive()
		fs.Logger.LogCompareReceive(msg, shadowMsg, shadowErr)
	}()

	return msg, nil
}
