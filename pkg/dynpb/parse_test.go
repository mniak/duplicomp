package dynpb

import (
	_ "embed"
	"fmt"
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
		// int32
		{
			Index: 1,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: 79,
			},
		},
		// String
		{
			Index: 2,
			ProtoValue: ProtoValue{
				Type:  TypeBytes,
				Bytes: []byte("Howdy, planet!"),
			},
		},
		// Booleans
		{
			Index: 3,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 1, // true
			},
		},
		{
			Index: 4,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 0, // false
			},
		},
		// Enum
		{
			Index: 5,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 2, // YELLOW
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
				Varint: EncodeTwosComplement(-42),
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
				Varint: EncodeTwosComplement(-1234567890123456789),
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
				Varint: EncodeZigZag(12345),
			},
		},
		{
			Index: 8,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: EncodeZigZag(-12345),
			},
		},
		{
			Index: 9,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: EncodeZigZag(98765432109876543),
			},
		},
		{
			Index: 10,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: EncodeZigZag(-98765432109876543),
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
				Fixed32: uint32(EncodeTwosComplement(123456789)),
			},
		},
		{
			Index: 14,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: uint32(EncodeTwosComplement(-123456789)),
			},
		},
		{
			Index: 15,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: EncodeTwosComplement(987654321012345678),
			},
		},
		{
			Index: 16,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: EncodeTwosComplement(-987654321012345678),
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
				Fixed32: EncodeFloat(3.1415926),
			},
		},
		{
			Index: 2,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: EncodeFloat(-3.1415926),
			},
		},
		// double
		{
			Index: 3,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: EncodeDouble(1.6180339887498),
			},
		},
		{
			Index: 4,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: EncodeDouble(-1.6180339887498),
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

func TestParseToMapWithHints_Example_Basic(t *testing.T) {
	ex := LoadExample("Basic")

	t.Run("No hints", func(t *testing.T) {
		parsed, err := parseToMapWithHints(ex.Bytes, make(HintMap))
		require.NoError(t, err)

		expected := map[int]any{
			// int32
			1: uint32(79),
			// String
			2: []byte("Howdy, planet!"),
			// Booleans
			3: uint64(1), // true
			4: uint64(0), // false
			// Enum
			5: uint64(2), // YELLOW
		}
		assert.Equal(t, expected, parsed)
	})

	t.Run("All hints", func(t *testing.T) {
		parsed, err := parseToMapWithHints(ex.Bytes, HintMap{
			1: HintInt{},
			2: HintString{},
			3: HintBoolean{},
			4: HintBoolean{},
			5: HintInt{},
		})
		require.NoError(t, err)

		expected := map[int]any{
			// int32
			1: int(79),
			// String
			2: string("Howdy, planet!"),
			// Booleans
			3: true,  // true
			4: false, // false
			// Enum
			5: uint64(2), // YELLOW
		}
		assert.Equal(t, expected, parsed)
	})
}
