package duplicomp_test

import (
	"io"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp"
	"github.com/mniak/duplicomp/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestForwarder2_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInStream := mocks.NewMockStream(ctrl)

	fakeMessage1 := new(duplicomp.Message)
	gofakeit.Struct(fakeMessage1)
	mockInStream.EXPECT().
		Receive().
		Return(fakeMessage1, nil)

	mockInStream.EXPECT().
		Receive().
		Return(nil, io.EOF)

	fw := duplicomp.Forwarder2{
		InboundStream: mockInStream,
	}
	err := fw.Run()
	require.NoError(t, err)
}
