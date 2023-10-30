package dynpb

import "errors"

type TypeHint interface {
	Apply(value any) (any, error)
}

type (
	HintInt           struct{}
	HintFloat         struct{}
	HintStruct        struct{}
	HintString        struct{}
	HintObject[T any] struct{}
	HintBoolean       struct{}
)

func (h HintInt) Apply(value any) (any, error) {
	var result int
	switch v := value.(type) {
	case int32:
		result = int(v)
	case uint32:
		result = int(v)
	case int64:
		result = int(v)
	case uint64:
		result = int(v)
	default:
		return value, errors.New("could not appy hint: Int")
	}
	return result, nil
}

func (h HintFloat) Apply(value any) (any, error) {
	return value, nil
}

func (h HintStruct) Apply(value any) (any, error) {
	return value, nil
}

func (h HintString) Apply(value any) (any, error) {
	var result string
	switch v := value.(type) {
	case []byte:
		result = string(v)
	default:
		return value, errors.New("could not appy hint: String")
	}
	return result, nil
}

func (h HintObject[T]) Apply(value any) (any, error) {
	return value, nil
}

func (h HintBoolean) Apply(value any) (any, error) {
	return value, nil
}
