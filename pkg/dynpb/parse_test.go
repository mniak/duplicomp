package dynpb

import (
	_ "embed"
	"fmt"
	"math"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestBuildDescriptor(t *testing.T) {
	desc := BuildDescriptor()
	assert.NotNil(t, desc)
	assert.NotEmpty(t, desc)
}

func BuildDescriptor() protoreflect.MessageDescriptor {
	return nil
}

func TestParseProto_Example_Basic(t *testing.T) {
	ex := LoadExample("Basic")
	parsed, err := parseProtoBytes(ex.Bytes)
	require.NoError(t, err)

	expected := []IndexedProtoValue{
		{
			Index: 1,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: 79,
			},
		},
		{
			Index: 2,
			ProtoValue: ProtoValue{
				Type:  TypeBytes,
				Bytes: []byte("Howdy, planet!"),
			},
		},
	}
	assert.Equal(t, expected, parsed)
}

func TestParseProto_Example_Integers(t *testing.T) {
	ex := LoadExample("Integers")
	parsed, err := parseProtoBytes(ex.Bytes)
	require.NoError(t, err)

	// Pay attention because each integer type represent negatives differently.
	// intN and sfixedN uses two's-complement encoding
	// sintN and sfixedN uses zigzag encoding
	// https://protobuf.dev/programming-guides/encoding/#signed-ints
	expected := []IndexedProtoValue{
		// intN uses two's-compement for negative numbers
		{
			Index: 1,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 42,
			},
		},
		{
			Index: 2,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: twosComplement(-42),
			},
		},
		{
			Index: 3,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 1234567890123456789,
			},
		},
		{
			Index: 4,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: twosComplement(-1234567890123456789),
			},
		},
		// uintN does not use negative, so they dont need encoding
		{
			Index: 5,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 12345,
			},
		},
		{
			Index: 6,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 98765432109876543,
			},
		},
		// sintN uses zig zag for negative numbers
		{
			Index: 7,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: zigzag(12345),
			},
		},
		{
			Index: 8,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: zigzag(-12345),
			},
		},
		{
			Index: 9,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: zigzag(98765432109876543),
			},
		},
		{
			Index: 10,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: zigzag(-98765432109876543),
			},
		},
		// fixedN does not have negative numbers
		{
			Index: 11,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: 123456789,
			},
		},
		{
			Index: 12,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: 987654321012345678,
			},
		},
		// sfixedN uses two's complement for negative numbers
		{
			Index: 13,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: uint32(twosComplement(123456789)),
			},
		},
		{
			Index: 14,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: uint32(twosComplement(-123456789)),
			},
		},
		{
			Index: 15,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: twosComplement(987654321012345678),
			},
		},
		{
			Index: 16,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: twosComplement(-987654321012345678),
			},
		},
	}

	expectedMap := lo.SliceToMap[IndexedProtoValue, int](expected, func(v IndexedProtoValue) (int, ProtoValue) {
		return v.Index, v.ProtoValue
	})
	parsedMap := lo.SliceToMap[IndexedProtoValue, int](parsed, func(v IndexedProtoValue) (int, ProtoValue) {
		return v.Index, v.ProtoValue
	})
	for fieldnum, expval := range expectedMap {
		t.Run(fmt.Sprintf("Field %d", fieldnum), func(t *testing.T) {
			actual := parsedMap[fieldnum]
			assert.Equal(t, expval, actual)
		})
	}

	assert.Equal(t, expected, parsed)
}

func float(f float32) uint32 {
	b := math.Float32bits(f)
	return b
}

func double(f float64) uint64 {
	b := math.Float64bits(f)
	return b
}

func TestParseProto_Example_Floats(t *testing.T) {
	ex := LoadExample("Floats")
	parsed, err := parseProtoBytes(ex.Bytes)
	require.NoError(t, err)

	expected := []IndexedProtoValue{
		// float
		{
			Index: 1,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: float(3.1415926),
			},
		},
		{
			Index: 2,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: float(-3.1415926),
			},
		},
		// double
		{
			Index: 3,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: double(1.6180339887498),
			},
		},
		{
			Index: 4,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: double(-1.6180339887498),
			},
		},
	}

	expectedMap := lo.SliceToMap[IndexedProtoValue, int](expected, func(v IndexedProtoValue) (int, ProtoValue) {
		return v.Index, v.ProtoValue
	})
	parsedMap := lo.SliceToMap[IndexedProtoValue, int](parsed, func(v IndexedProtoValue) (int, ProtoValue) {
		return v.Index, v.ProtoValue
	})
	for fieldnum, expval := range expectedMap {
		t.Run(fmt.Sprintf("Field %d", fieldnum), func(t *testing.T) {
			actual := parsedMap[fieldnum]
			assert.Equal(t, expval, actual)
		})
	}

	assert.Equal(t, expected, parsed)
}

func zigzag(v int64) uint64 {
	if v >= 0 {
		return uint64(v * 2)
	} else {
		return uint64(v*-2 - 1)
	}
}

func twosComplement(v int64) uint64 {
	return uint64(v)
}
