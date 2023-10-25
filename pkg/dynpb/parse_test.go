package dynpb

import (
	_ "embed"
	"encoding/base64"
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

func s64(n int64) uint64 {
	return uint64(n)
}

func u64(n uint64) uint64 {
	return n
}

func f64(n float64) uint64 {
	return uint64(n)
}

func s32(n int32) uint32 {
	return uint32(n)
}

func u32(n uint32) uint32 {
	return n
}

func f32(n float32) uint64 {
	return uint64(n)
}

func TestParseProto_Example1(t *testing.T) {
	parsed, err := parseProtoBytes(example1.Bytes)
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

func TestParseProto_Example2(t *testing.T) {
	parsed, err := parseProtoBytes(example2.Bytes)
	require.NoError(t, err)

	expected := []IndexedProtoValue{
		{
			Index: 1,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: s64(42),
			},
		},
		{
			Index: 2,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: s64(1234567890123456789),
			},
		},
		{
			Index: 3,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: u64(12345),
			},
		},
		{
			Index: 4,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: u64(98765432109876543),
			},
		},
		{
			Index: 5,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: s64(-12345),
			},
		},
		{
			Index: 6,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: s64(-98765432109876543),
			},
		},
		{
			Index: 7,
			ProtoValue: ProtoValue{
				Type:    TypeVarint,
				Fixed32: u32(123456789),
			},
		},
		{
			Index: 8,
			ProtoValue: ProtoValue{
				Type:    TypeVarint,
				Fixed64: u64(987654321012345678),
			},
		},
		{
			Index: 9,
			ProtoValue: ProtoValue{
				Type:    TypeFixed32,
				Fixed32: s32(-123456789),
			},
		},
		{
			Index: 10,
			ProtoValue: ProtoValue{
				Type:    TypeFixed64,
				Fixed64: s64(-987654321012345678),
			},
		},
		{
			Index: 11,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: f32(3.14159),
			},
		},
		{
			Index: 12,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: f64(2.71828),
			},
		},
		{
			Index: 13,
			ProtoValue: ProtoValue{
				Type:   TypeVarint,
				Varint: 1,
			},
		},
		{
			Index: 14,
			ProtoValue: ProtoValue{
				Type:  TypeBytes,
				Bytes: []byte("Hello, world!"),
			},
		},
		{
			Index: 15,
			ProtoValue: ProtoValue{
				Type:  TypeBytes,
				Bytes: lo.Must(base64.StdEncoding.DecodeString("SGVsbG8sIHdvcmxkIQ==")),
			},
		},

		// repeated int32 field_int32_list = 16;
		// repeated int64 field_int64_list = 17;
		// repeated uint32 field_uint32_list = 18;
		// repeated uint64 field_uint64_list = 19;
		// repeated sint32 field_sint32_list = 20;
		// repeated sint64 field_sint64_list = 21;
		// repeated fixed32 field_fixed32_list = 22;
		// repeated fixed64 field_fixed64_list = 23;
		// repeated sfixed32 field_sfixed32_list = 24;
		// repeated sfixed64 field_sfixed64_list = 25;
		// repeated float field_float_list = 26;
		// repeated double field_double_list = 27;
		// repeated bool field_bool_list = 28;
		// repeated string field_string_list = 29;
		// repeated bytes field_bytes_list = 30;

		// optional Object field_object = 31;
		// repeated Object field_object_list = 32;
	}
	assert.Equal(t, expected, parsed)
}
