package duplicomp

import "google.golang.org/protobuf/proto"

//go:generate mockgen -package=duplicomp -destination=mock_comparator_test.go . Comparator
type Comparator interface {
	Compare(
		msg1 proto.Message, err1 error,
		msg2 proto.Message, err2 error,
	)
}
