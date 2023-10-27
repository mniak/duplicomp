package dynpb

type TypeHint interface {
	Apply(value any) (any, error)
}

type (
	HintInt           struct{}
	HintFloat         struct{}
	HintStruct        struct{}
	HintString        struct{}
	HintObject[T any] struct{}
)

func (h HintInt) Apply(value any) (any, error) {
	return value, nil
}

func (h HintFloat) Apply(value any) (any, error) {
	return value, nil
}

func (h HintStruct) Apply(value any) (any, error) {
	return value, nil
}

func (h HintString) Apply(value any) (any, error) {
	return value, nil
}

func (h HintObject[T]) Apply(value any) (any, error) {
	return value, nil
}
