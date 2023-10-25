package duplicomp

//go:generate mockgen -package=duplicomp -destination=mock_comparator_test.go . Comparator
type Comparator interface {
	Compare(
		data1 []byte, err1 error,
		data2 []byte, err2 error,
	) error
}
