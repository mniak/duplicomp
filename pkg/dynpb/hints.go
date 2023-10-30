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
	HintBool          struct{}
	HintEnum[T ~int]  struct{}
)

func (h HintInt) Apply(value any) (any, error) {
	switch v := value.(type) {
	case int32:
		return int(v), nil
	case uint32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint64:
		return int(v), nil
	default:
		return value, errors.New("could not appy hint: Int")
	}
}

func (h HintFloat) Apply(value any) (any, error) {
	return value, nil
}

func (h HintStruct) Apply(value any) (any, error) {
	return value, nil
}

func (h HintString) Apply(value any) (any, error) {
	switch v := value.(type) {
	case []byte:
		return string(v), nil
	default:
		return value, errors.New("could not appy hint: String")
	}
}

func (h HintObject[T]) Apply(value any) (any, error) {
	return value, nil
}

func (h HintBool) Apply(value any) (any, error) {
	switch v := value.(type) {
	case int32:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	default:
		return value, errors.New("could not appy hint: Bool")
	}
}

func (h HintEnum[T]) Apply(value any) (any, error) {
	switch v := value.(type) {
	case int32:
		return T(v), nil
	case uint32:
		return T(v), nil
	case int64:
		return T(v), nil
	case uint64:
		return T(v), nil
	default:
		return value, errors.New("could not appy hint: Enum")
	}
}
