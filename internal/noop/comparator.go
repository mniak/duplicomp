package noop

func Comparator() _NoOpComparator {
	return _NoOpComparator{}
}

type _NoOpComparator struct{}

func (_NoOpComparator) Compare(
	primaryMsg []byte, primaryError error,
	shadowMsg []byte, shadowError error,
) error {
	return nil
}
