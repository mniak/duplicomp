package ps121

//go:generate mockgen -package=ps121 -destination=mock_comparator_test.go . Comparator
type Comparator interface {
	Compare(
		methodName string,
		primaryData []byte, primaryError error,
		shadowData []byte, shadowError error,
	) error
}
