package duplicomp

import (
	"google.golang.org/protobuf/proto"
)

//go:generate mockgen -package=duplicomp -destination=mock_shadowlogger_test.go . ShadowLogger
type ShadowLogger interface {
	LogSendError(error)
}

type ShadowStream struct {
	primaryStream Stream
	shadowStream  Stream
	shadowLogger  ShadowLogger
}

func (fs *ShadowStream) Send(m proto.Message) error {
	// var wg sync.WaitGroup
	// wg.Add(2)

	// var err error
	// go func() {
	// 	err = fs.primaryStream.Send(m)
	// 	wg.Done()
	// }()

	// var shadowError error
	// go func() {
	// 	shadowError = fs.shadowStream.Send(m)
	// 	wg.Done()
	// }()
	// _ = shadowError
	// wg.Wait()
	// if shadowError != nil {
	// 	fs.shadowLogger.LogSendError(shadowError)
	// }

	err := fs.primaryStream.Send(m)
	if err != nil {
		return err
	}

	go func() {
		err = fs.shadowStream.Send(m)
		// if err != nil {
		// 	return err
		// }
	}()

	return err
}

func (fs *ShadowStream) Receive() (proto.Message, error) {
	return nil, nil
}
