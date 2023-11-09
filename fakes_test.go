package duplicomp

import (
	"github.com/brianvoe/gofakeit/v6"
	pbany "github.com/golang/protobuf/ptypes/any"
)

func NewFakeProtoMessage(bytes []byte) any {
	msg := new(pbany.Any)
	gofakeit.Struct(msg)
	msg.ProtoReflect().SetUnknown(bytes)
	return msg
}
