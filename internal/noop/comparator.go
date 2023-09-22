package noop

func Comparator() _NoOpComparator {
	return _NoOpComparator{}
}

type _NoOpComparator struct{}

func (_NoOpComparator) Compare(
	msg1 []byte, err1 error,
	msg2 []byte, err2 error,
) error {
	return nil
}
