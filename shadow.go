package duplicomp

import (
	"google.golang.org/protobuf/proto"
)

//go:generate mockgen -package=duplicomp -destination=mock_shadowlogger_test.go . ShadowLogger
type ShadowLogger interface {
	LogSendFailure(error)
}

type StreamWithShadow struct {
	inner  Stream
	shadow Stream
	logger ShadowLogger
}

func (fs *StreamWithShadow) Send(m proto.Message) error {
	// var wg sync.WaitGroup
	// wg.Add(2)

	// var err error
	// go func() {
	// 	err = fs.innerStream.Send(m)
	// 	wg.Done()
	// }()

	// var shadowError error
	// go func() {
	// 	shadowError = fs.shadow.Send(m)
	// 	wg.Done()
	// }()
	// _ = shadowError
	// wg.Wait()
	// if shadowError != nil {
	// 	fs.shadowLogger.LogSendError(shadowError)
	// }

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
	return nil, nil
}
