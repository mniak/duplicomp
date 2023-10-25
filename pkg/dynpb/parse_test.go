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

func TestParseProto_Example1(t *testing.T) {
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

func zigzag(v int64) uint64 {
	if v >= 0 {
		return uint64(v * 2)
	} else {
		return uint64(v*-2 - 1)
	}
}

func onescomp(v int64) uint64 {
	return uint64(v)
}

// func sint32(v int) int32 {
// 	return int32(onescomp(v))
// }

func TestParseProto_Example2(t *testing.T) {
	ex := LoadExample("Integers")
	parsed, err := parseProtoBytes(ex.Bytes)
	require.NoError(t, err)

	expected := []IndexedProtoValue{
		// intN uses ones-compement for negative numbers
		{
			Index: 1,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 42,
			},
		},
		{
			Index: 101,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: onescomp(-42),
			},
		},
		{
			Index: 2,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 1234567890123456789,
			},
		},
		{
			Index: 102,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: onescomp(-1234567890123456789),
			},
		},
		// {
		// 	Index: 3,
		// 	ProtoValue: ProtoValue{
		// 		Type:   TypeVarint,
		// 		Varint: zigzag(12345),
		// 	},
		// },
		// {
		// 	Index: 4,
		// 	ProtoValue: ProtoValue{
		// 		Type:   TypeVarint,
		// 		Varint: zigzag(98765432109876543),
		// 	},
		// },
		// {
		// 	Index: 5,
		// 	ProtoValue: ProtoValue{
		// 		Type:   TypeVarint,
		// 		Varint: zigzag(-12345),
		// 	},
		// },
		// {
		// 	Index: 6,
		// 	ProtoValue: ProtoValue{
		// 		Type:   TypeVarint,
		// 		Varint: zigzag(-98765432109876543),
		// 	},
		// },
		// {
		// 	Index: 7,
		// 	ProtoValue: ProtoValue{
		// 		Type:    TypeFixed32,
		// 		Fixed32: 123456789,
		// 	},
		// },
		// {
		// 	Index: 8,
		// 	ProtoValue: ProtoValue{
		// 		Type:    TypeFixed64,
		// 		Fixed64: 987654321012345678,
		// 	},
		// },
		{
			Index: 9,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: uint32(onescomp(-123456789)),
			},
		},
		// {
		// 	Index: 10,
		// 	ProtoValue: ProtoValue{
		// 		Type:    TypeFixed64,
		// 		Fixed64: zigzag(-987654321012345678),
		// 	},
		// },
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

	// assert.Equal(t, expected, parsed)
}
