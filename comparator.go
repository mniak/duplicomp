package duplicomp

import "google.golang.org/protobuf/proto"

//go:generate mockgen -package=duplicomp -destination=mock_shadowlogger_test.go . ShadowLogger
type ShadowLogger interface {
	LogSendFailure(error)
	LogCompareReceive(primaryMsg, shadowMsg proto.Message, shadowErr error)
}
