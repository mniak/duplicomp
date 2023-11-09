package duplicomp

//go:generate mockgen -package=duplicomp -destination=mock_comparator_test.go . Comparator
type Comparator interface {
	Compare(
		primaryData []byte, primaryError error,
		shadowData []byte, shadowError error,
	) error
}
