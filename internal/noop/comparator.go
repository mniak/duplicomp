package noop

import "google.golang.org/protobuf/proto"

func Comparator() _NoOpComparator {
	return _NoOpComparator{}
}

type _NoOpComparator struct{}

func (_NoOpComparator) Compare(
	msg1 proto.Message, err1 error,
	msg2 proto.Message, err2 error,
) {
}
