package dynpb

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

type ProtoType string

const (
	TypeBytes   ProtoType = "len"
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
	Group   ProtoMap
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

type invalidType struct{}

func (v ProtoValue) RawValue() any {
	switch v.Type {
	case TypeBytes:
		return v.Bytes
	case TypeVarint:
		return v.Varint
	case TypeFixed32:
		return v.Fixed32
	case TypeFixed64:
		return v.Fixed64
	case TypeGroup:
		return v.Group.ProtoMapToMap()

	default:
		return invalidType{}
	}
}

type (
	ProtoMap          []IndexedProtoValue
	Object            = map[int]any
	IndexedProtoValue struct {
		Index int
		ProtoValue
	}
)

func (pm ProtoMap) ProtoMapToMap() Object {
	return lo.Associate[IndexedProtoValue, int, any](pm, func(item IndexedProtoValue) (int, any) {
		return item.Index, item.RawValue()
	})
}

func ParseProtoMessage(m proto.Message) (ProtoMap, error) {
	unknownBytes := m.ProtoReflect().GetUnknown()
	fields, err := parseProtoBytes(unknownBytes)
	return fields, err
}

func parseProtoBytes(b []byte) (ProtoMap, error) {
	var result ProtoMap
	const dec = 10
	const hex = 16
	for len(b) > 0 {
		num, wtype, length := protowire.ConsumeTag(b)
		if length < 0 {
			return nil, errors.New("failed to consume tag")
		}
		b = b[length:]

		switch wtype {
		case protowire.VarintType:
			var v uint64
			v, length = protowire.ConsumeVarint(b)
			if length < 0 {
				return nil, fmt.Errorf("failed to parse varint %d: %s", num, protowire.ParseError(length))
			}
			result = append(result, IndexedProtoValue{
				Index: int(num),
				ProtoValue: ProtoValue{
					Type:   TypeVarint,
					Varint: v,
				},
			})
		case protowire.Fixed32Type:
			var v uint32
			v, length = protowire.ConsumeFixed32(b)
			if length < 0 {
				return nil, fmt.Errorf("failed to parse fixed32 %d: %s", num, protowire.ParseError(length))
			}
			result = append(result, IndexedProtoValue{
				Index: int(num),
				ProtoValue: ProtoValue{
					Type:    TypeFixed32,
					Fixed32: v,
				},
			})
		case protowire.Fixed64Type:
			var v uint64
			v, length = protowire.ConsumeFixed64(b)
			if length < 0 {
				return nil, fmt.Errorf("failed to parse fixed64 %d: %s", num, protowire.ParseError(length))
			}
			result = append(result, IndexedProtoValue{
				Index: int(num),
				ProtoValue: ProtoValue{
					Type:    TypeFixed64,
					Fixed64: v,
				},
			})
		case protowire.BytesType:
			var v []byte
			v, length = protowire.ConsumeBytes(b)
			if length < 0 {
				return nil, fmt.Errorf("failed to parse bytes %d: %s", num, protowire.ParseError(length))
			}
			result = append(result, IndexedProtoValue{
				Index: int(num),
				ProtoValue: ProtoValue{
					Type:  TypeBytes,
					Bytes: v,
				},
			})
		case protowire.StartGroupType:
			var g []byte
			g, length = protowire.ConsumeGroup(num, b)
			if length < 0 {
				return nil, fmt.Errorf("failed to parse group %d: %s", num, protowire.ParseError(length))
			}
			v, err := parseProtoBytes(g)
			if err != nil {
				return nil, err
			}
			result = append(result, IndexedProtoValue{
				Index: int(num),
				ProtoValue: ProtoValue{
					Type:  TypeGroup,
					Group: v,
				},
			})
		default:
			return nil, fmt.Errorf("error parsing unknown field wire type: %v", wtype)
		}

		b = b[length:]
	}
	return result, nil
}

type (
	ProtoHintMap  = map[int]ProtoDataHint
	ProtoDataHint struct {
		Name      string
		SubFields map[int]ProtoDataHint
	}
)

func parseToMapWithHints(data []byte, hints HintMap) (Object, error) {
	if hints == nil {
		hints = make(HintMap)
	}
	protoMap, err := parseProtoBytes(data)
	if err != nil {
		return nil, err
	}
	result := make(Object)
	for _, field := range protoMap {
		value := field.RawValue()

		if hint, hasHint := hints[field.Index]; hasHint {
			var err error
			current := result[field.Index]
			value, err = hint.Apply(current, field.RawValue())
			if err != nil {
				return nil, err
			}
		}

		result[field.Index] = value
	}
	return result, nil
}
