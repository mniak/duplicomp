package dynpb

import (
	"fmt"

	"github.com/pkg/errors"
)

type IntEncoding uint8

const (
	Unsigned       IntEncoding = 0
	TwosComplement IntEncoding = 1
	ZigZag         IntEncoding = 2
)

type IntegerHint string

const (
	HintInt32       IntegerHint = "Int32"
	HintInt32ZigZag IntegerHint = "Int32ZigZag"
	HintUInt32      IntegerHint = "UInt32"

	HintInt64       IntegerHint = "Int64"
	HintInt64ZigZag IntegerHint = "Int64ZigZag"
	HintUInt64      IntegerHint = "UInt64"

	// HintFloat IntegerHint = "Float"
	// HintString SimpleHint = "String"
	HintBool IntegerHint = "Bool"
)

func (h IntegerHint) getValue(value any) (uint64, error) {
	switch v := value.(type) {
	case uint32:
		return uint64(v), nil
	case uint64:
		return uint64(v), nil
	default:
		return 0, errors.New("could get integer value for hint")
	}
}

func (h IntegerHint) Apply(value any) (any, error) {
	val, err := h.getValue(value)
	if err != nil {
		return nil, err
	}
	switch h {
	case HintInt32:
		return int32(val), nil
	case HintInt32ZigZag:
		return int32(DecodeZigZag(val)), nil
	case HintUInt32:
		return uint32(val), nil

	case HintInt64:
		return int64(val), nil
	case HintInt64ZigZag:
		return int64(DecodeZigZag(val)), nil
	case HintUInt64:
		return uint64(val), nil

	// case HintFloat:
	// 	return float(val), nil
	case HintBool:
		return val != 0, nil

	default:
		return nil, fmt.Errorf("invalid simple hint: %q", string(h))
	}
}

type (
	TypeHint interface {
		Apply(value any) (any, error)
	}
	HintEnum[T ~int] struct{}
)

type ByteSliceHint string

const (
	HintBytes  ByteSliceHint = "bytes"
	HintString ByteSliceHint = "string"
)

func (h ByteSliceHint) Apply(value any) (any, error) {
	bytes, ok := value.([]byte)
	if !ok {
		return nil, errors.New("could get byte slice value for hint")
	}

	switch h {
	case HintBytes:
		return bytes, nil
	case HintString:
		return string(bytes), nil
	default:
		return nil, fmt.Errorf("invalid byte slice hint: %q", string(h))
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
