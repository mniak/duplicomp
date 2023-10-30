package dynpb

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProto_Example_Basic(t *testing.T) {
	ex := LoadExample("Basic")
	parsed, err := parseProtoBytes(ex.Bytes)
	require.NoError(t, err)

	expected := ProtoMap{
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
	expected := ProtoMap{
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

	assert.Equal(t, expected, parsed)
}

func TestParseProto_Example_Floats(t *testing.T) {
	ex := LoadExample("Floats")
	parsed, err := parseProtoBytes(ex.Bytes)
	require.NoError(t, err)

	expected := ProtoMap{
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
		type Color1 int
		const YELLOW Color1 = 2
		const BLUE Color1 = 1

		parsed, err := parseToMapWithHints(ex.Bytes, HintMap{
			1: HintInt{},
			2: HintString{},
			3: HintBool{},
			4: HintBool{},
			5: HintEnum[Color1]{},
		})
		require.NoError(t, err)

		expected := map[int]any{
			// int
			1: int(79),
			// String
			2: string("Howdy, planet!"),
			// Booleans
			3: true,  // true
			4: false, // false
			// Enum
			5: YELLOW, // YELLOW
		}
		assert.Equal(t, expected, parsed)
	})
}

func TestParseToMapWithHints_Example_Integers(t *testing.T) {
	ex := LoadExample("Integers")
	parsed, err := parseToMapWithHints(
		ex.Bytes,
		HintMap{
			// intN uses two's-compement for negative numbers
			1: HintInt{Encoding: TwosComplement},
			2: HintInt{Encoding: TwosComplement},
			3: HintInt{Encoding: TwosComplement},
			4: HintInt{Encoding: TwosComplement},
			// uintN does not use negative, so they dont need encoding
			5: HintInt{},
			6: HintInt{},
			// sintN uses zig zag for negative numbers
			7:  HintInt{Encoding: ZigZag},
			8:  HintInt{Encoding: ZigZag},
			9:  HintInt{Encoding: ZigZag},
			10: HintInt{Encoding: ZigZag},
			// fixedN does not have negative numbers
			11: HintInt{},
			12: HintInt{},
			// sfixedN uses two's complement for negative numbers
			13: HintInt{Encoding: TwosComplement},
			14: HintInt{Encoding: TwosComplement},
			15: HintInt{Encoding: TwosComplement},
			16: HintInt{Encoding: TwosComplement},
		},
	)
	require.NoError(t, err)

	// Pay attention because each integer type represent negatives differently.
	// intN and sfixedN uses two's-complement encoding
	// sintN and sfixedN uses zigzag encoding
	// https://protobuf.dev/programming-guides/encoding/#signed-ints
	expected := map[int]any{
		// intN uses two's-compement for negative numbers
		1: 42,
		2: -42,
		3: 1234567890123456789,
		4: -1234567890123456789,
		// uintN does not use negative, so they dont need encoding
		5: 12345,
		6: 98765432109876543,
		// sintN uses zig zag for negative numbers
		7:  12345,
		8:  -12345,
		9:  98765432109876543,
		10: -98765432109876543,
		// fixedN does not have negative numbers
		11: 123456789,
		12: 987654321012345678,
		// sfixedN uses two's complement for negative numbers
		13: 123456789,
		14: -123456789,
		15: 987654321012345678,
		16: -987654321012345678,
	}

	assert.Equal(t, expected, parsed)
}
