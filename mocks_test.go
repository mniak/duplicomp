package duplicomp

import "google.golang.org/protobuf/reflect/protoreflect"

//go:generate mockgen -package=duplicomp -destination=mock_proto_message_test.go "google.golang.org/protobuf/proto" Message
//go:generate mockgen -package=duplicomp -destination=mock_protoreflect_message_test.go "google.golang.org/protobuf/reflect/protoreflect" Message

// func NewFakeMessage(ctrl *gomock.Controller, payload []byte) *MockMessage {
// 	mockProtoMessage := NewMockProtoMessage(ctrl)
// 	mockProtoMessage.EXPECT().ProtoReflect()

// 	mockMessage := NewMockMessage(ctrl)
// 	mockMessage.EXPECT().ProtoReflect().AnyTimes().Return(mockProtoMessage)
// }

type FakeProtoreflectMessage struct {
	protoreflect.Message
}

// func (self FakeProtoreflectMessage) GetUnknown() protoreflect.RawFields {
// 	return self.
// }
