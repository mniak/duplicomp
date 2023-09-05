package duplicomp

import (
	"io"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPipe_HappyPath(t *testing.T) {
	t.Run("With two messages", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeMessage1 := new(any.Any)
		gofakeit.Struct(*fakeMessage1)
		require.NotNil(t, fakeMessage1)

		fakeMessage2 := new(any.Any)
		gofakeit.Struct(&fakeMessage2)
		require.NotNil(t, fakeMessage2)

		require.NotEqual(t, fakeMessage1, fakeMessage2)
		require.NotEqual(t, &fakeMessage1, &fakeMessage2)

		inMock := NewMockStream(ctrl)
		inMock.EXPECT().
			Receive().
			Return(fakeMessage1, nil)
		inMock.EXPECT().
			Receive().
			Return(nil, io.EOF)

		outMock := NewMockStream(ctrl)
		outMock.EXPECT().Send(fakeMessage1).Return(nil)

		fw := Pipe{
			In:  inMock,
			Out: outMock,
		}
		err := fw.Run()
		require.NoError(t, err)
	})
}
