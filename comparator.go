package duplicomp

//go:generate mockgen -package=duplicomp -destination=mock_comparator_test.go . Comparator
type Comparator interface {
	Compare(
		msg1 []byte, err1 error,
		msg2 []byte, err2 error,
	) error
}
