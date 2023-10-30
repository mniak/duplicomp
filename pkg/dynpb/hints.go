package dynpb

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	HintMap  map[int]TypeHint
	TypeHint interface {
		Apply(value any) (any, error)
	}
)

func (h HintMap) Apply(value any) (any, error) {
	bytes, ok := value.([]byte)
	if !ok {
		return nil, errors.New("could get byte slice value for hint")
	}
	return parseToMapWithHints(bytes, h)
}

type IntEncoding uint8

const (
	Unsigned       IntEncoding = 0
	TwosComplement IntEncoding = 1
	ZigZag         IntEncoding = 2
)

type NumericHint string

const (
	HintInt32       NumericHint = "Int32"
	HintInt32ZigZag NumericHint = "Int32ZigZag"
	HintUInt32      NumericHint = "UInt32"

	HintInt64       NumericHint = "Int64"
	HintInt64ZigZag NumericHint = "Int64ZigZag"
	HintUInt64      NumericHint = "UInt64"

	HintFloat  NumericHint = "Float"
	HintDouble NumericHint = "Double"

	HintBool NumericHint = "Bool"
)

func (h NumericHint) getValue(value any) (uint64, error) {
	switch v := value.(type) {
	case uint32:
		return uint64(v), nil
	case uint64:
		return uint64(v), nil
	default:
		return 0, errors.New("could get number value for hint")
	}
}

func (h NumericHint) Apply(value any) (any, error) {
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

	case HintFloat:
		return DecodeFloat(val), nil
	case HintDouble:
		return DecodeDouble(val), nil

	case HintBool:
		return val != 0, nil

	default:
		return nil, fmt.Errorf("invalid numeric hint: %q", string(h))
	}
}

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

type HintEnum[T ~int] struct{}

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
