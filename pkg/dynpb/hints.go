package dynpb

import "errors"

type IntEncoding uint8

const (
	Unsigned       IntEncoding = 0
	TwosComplement IntEncoding = 1
	ZigZag         IntEncoding = 2
)

type (
	TypeHint interface {
		Apply(value any) (any, error)
	}
	HintInt struct {
		Encoding IntEncoding
	}
	hintIntZigZag         struct{}
	hintIntTwosComplement struct{}
	hintIntUnsigned       struct{}

	HintFloat        struct{}
	HintString       struct{}
	HintBool         struct{}
	HintEnum[T ~int] struct{}
)

func (h hintIntZigZag) Apply(value any) (any, error) {
	switch v := value.(type) {
	case int32:
		return int(DecodeZigZag(uint64(v))), nil
	case uint32:
		return int(DecodeZigZag(uint64(v))), nil
	case int64:
		return int(DecodeZigZag(uint64(v))), nil
	case uint64:
		return int(DecodeZigZag(uint64(v))), nil
	default:
		return value, errors.New("could not appy hint: Int-zigzag")
	}
}

func (h hintIntTwosComplement) Apply(value any) (any, error) {
	switch v := value.(type) {
	case uint32:
		return int(int32(v)), nil
	case uint64:
		return int(int64(v)), nil
	default:
		return value, errors.New("could not appy hint: Int")
	}
}

func (h hintIntUnsigned) Apply(value any) (any, error) {
	switch v := value.(type) {
	case int32:
		return int(uint64(v)), nil
	case uint32:
		return int(uint64(v)), nil
	case int64:
		return int(uint64(v)), nil
	case uint64:
		return int(uint64(v)), nil
	default:
		return value, errors.New("could not appy hint: Int")
	}
}

func (h HintInt) Apply(value any) (any, error) {
	switch h.Encoding {
	case ZigZag:
		return hintIntZigZag{}.Apply(value)
	case TwosComplement:
		return hintIntTwosComplement{}.Apply(value)
	default:
		return hintIntUnsigned{}.Apply(value)
	}
}

func (h HintFloat) Apply(value any) (any, error) {
	return value, errors.New("could not appy hint: Float")
}

func (h HintString) Apply(value any) (any, error) {
	switch v := value.(type) {
	case []byte:
		return string(v), nil
	default:
		return value, errors.New("could not appy hint: String")
	}
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
