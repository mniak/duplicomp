package duplicomp

import (
	"fmt"
	"unicode/utf8"

	"github.com/elliotchance/orderedmap/v2"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

type ProtoType string

const (
	TypeBytes   ProtoType = "bytes"
	TypeVarint  ProtoType = "varint"
	TypeFixed32 ProtoType = "fixed32"
	TypeFixed64 ProtoType = "fixed64"
	TypeGroup   ProtoType = "group"
)

type ProtoValue struct {
	Type ProtoType

	Bytes   []byte
	Varint  uint64
	Fixed32 uint32
	Fixed64 uint64
	Group   *ProtoMap
}

func (v ProtoValue) String() string {
	switch v.Type {
	case TypeBytes:
		if utf8.Valid(v.Bytes) {
			return fmt.Sprintf("%q", string(v.Bytes))
		}
		return fmt.Sprintf("%2X", v.Bytes)
	case TypeVarint:
		return fmt.Sprint(v.Varint)
	case TypeFixed32:
		return fmt.Sprint(v.Fixed32)
	case TypeFixed64:
		return fmt.Sprint(v.Fixed64)
	case TypeGroup:
		return "<group>"

	default:
		return "<invalid type>"
	}
}

type ProtoMap = orderedmap.OrderedMap[int, ProtoValue]

func ParseProtoMessage(m proto.Message) (*ProtoMap, error) {
	unknownBytes := m.ProtoReflect().GetUnknown()
	fields, err := parseProtoBytes(unknownBytes)
	return fields, err
}

func parseProtoBytes(b []byte) (*ProtoMap, error) {
	result := orderedmap.NewOrderedMap[int, ProtoValue]()
	const dec = 10
	const hex = 16
	for len(b) > 0 {
		num, wtype, n := protowire.ConsumeTag(b)
		b = b[n:]

		switch wtype {
		case protowire.VarintType:
			var v uint64
			v, n = protowire.ConsumeVarint(b)
			result.Set(int(num), ProtoValue{
				Type:   TypeVarint,
				Varint: v,
			})
		case protowire.Fixed32Type:
			var v uint32
			v, n = protowire.ConsumeFixed32(b)
			result.Set(int(num), ProtoValue{
				Type:    TypeFixed32,
				Fixed32: v,
			})
		case protowire.Fixed64Type:
			var v uint64
			v, n = protowire.ConsumeFixed64(b)
			result.Set(int(num), ProtoValue{
				Type:    TypeFixed64,
				Fixed64: v,
			})
		case protowire.BytesType:
			var v []byte
			v, n = protowire.ConsumeBytes(b)
			result.Set(int(num), ProtoValue{
				Type:  TypeBytes,
				Bytes: v,
			})
		case protowire.StartGroupType:
			var g []byte
			g, n = protowire.ConsumeGroup(num, b)
			v, err := parseProtoBytes(g)
			if err != nil {
				return nil, err
			}
			result.Set(int(num), ProtoValue{
				Type:  TypeGroup,
				Group: v,
			})
		default:
			return nil, fmt.Errorf("error parsing unknown field wire type: %v", wtype)
		}

		b = b[n:]
	}
	return (*ProtoMap)(result), nil
}
