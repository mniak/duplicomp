package duplicomp

import (
	"io"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestForwarder2_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	var fakeMessage1 proto.Message
	gofakeit.Struct(&fakeMessage1)

	inMock := NewMockStream(ctrl)
	inMock.EXPECT().
		Receive().
		Return(fakeMessage1, nil)
	inMock.EXPECT().
		Receive().
		Return(nil, io.EOF)

	outMock1 := NewMockStream(ctrl)
	outMock1.EXPECT().Send(fakeMessage1).Return(nil)

	fw := Forwarder2{
		InboundStream:  inMock,
		OutboundStream: outMock1,
	}
	err := fw.Run()
	require.NoError(t, err)
}
