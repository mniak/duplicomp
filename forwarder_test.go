package duplicomp_test

import (
	"io"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp"
	"github.com/mniak/duplicomp/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestForwarder2_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	fakeMessage1 := new(proto.Message)
	gofakeit.Struct(fakeMessage1)

	inMock := mocks.NewMockStream(ctrl)
	inMock.EXPECT().
		Receive().
		Return(fakeMessage1, nil)
	inMock.EXPECT().
		Receive().
		Return(nil, io.EOF)

	outMock1 := mocks.NewMockStream(ctrl)
	outMock1.EXPECT().Send(fakeMessage1).Return(nil)

	fw := duplicomp.Forwarder2{
		InboundStream:  inMock,
		OutboundStream: outMock1,
	}
	err := fw.Run()
	require.NoError(t, err)
}
