package duplicomp

import (
	"google.golang.org/protobuf/proto"
)

//go:generate mockgen -package=duplicomp -destination=mock_shadowlogger_test.go . ShadowLogger
type ShadowLogger interface {
	LogSendFailure(error)
	LogCompareReceive(msg, shadowMsg proto.Message, shadowErr error)
}

type StreamWithShadow struct {
	inner  Stream
	shadow Stream
	logger ShadowLogger
}

func (fs *StreamWithShadow) Send(m proto.Message) error {
	err := fs.inner.Send(m)
	if err != nil {
		return err
	}

	go func() {
		err := fs.shadow.Send(m)
		if err != nil {
			fs.logger.LogSendFailure(err)
		}
	}()

	return nil
}

func (fs *StreamWithShadow) Receive() (proto.Message, error) {
	msg, err := fs.inner.Receive()
	if err != nil {
		return msg, err
	}

	go func() {
		shadowMsg, shadowErr := fs.shadow.Receive()
		fs.logger.LogCompareReceive(msg, shadowMsg, shadowErr)
	}()

	return msg, nil
}
