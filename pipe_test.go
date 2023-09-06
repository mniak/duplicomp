package duplicomp

import (
	"errors"
	"io"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPipe_HappyPath(t *testing.T) {
	t.Run("With no messages", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		inMock := NewMockStream(ctrl)
		inMock.EXPECT().
			Receive().
			Return(nil, io.EOF)

		outMock := NewMockStream(ctrl)

		fw := Pipe{
			In:  inMock,
			Out: outMock,
		}
		err := fw.Run()
		require.NoError(t, err)
	})
	t.Run("With one message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeMessage := new(any.Any)
		gofakeit.Struct(fakeMessage)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)

		inMock := NewMockStream(ctrl)
		inMock.EXPECT().
			Receive().
			Return(fakeMessage, nil)
		inMock.EXPECT().
			Receive().
			Return(nil, io.EOF)

		outMock := NewMockStream(ctrl)
		outMock.EXPECT().Send(fakeMessage).Return(nil)

		fw := Pipe{
			In:  inMock,
			Out: outMock,
		}
		err := fw.Run()
		require.NoError(t, err)
	})
	t.Run("With two messages", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeMessage1 := new(any.Any)
		gofakeit.Struct(fakeMessage1)
		require.NotNil(t, fakeMessage1)
		require.NotEmpty(t, fakeMessage1)

		fakeMessage2 := new(any.Any)
		gofakeit.Struct(fakeMessage2)
		require.NotNil(t, fakeMessage2)
		require.NotEmpty(t, fakeMessage2)

		require.NotEqual(t, fakeMessage1, fakeMessage2)
		require.NotEqual(t, &fakeMessage1, &fakeMessage2)

		inMock := NewMockStream(ctrl)
		inMock.EXPECT().
			Receive().
			Return(fakeMessage1, nil)
		inMock.EXPECT().
			Receive().
			Return(fakeMessage2, nil)
		inMock.EXPECT().
			Receive().
			Return(nil, io.EOF)

		outMock := NewMockStream(ctrl)
		outMock.EXPECT().Send(fakeMessage1).Return(nil)
		outMock.EXPECT().Send(fakeMessage2).Return(nil)

		fw := Pipe{
			In:  inMock,
			Out: outMock,
		}
		err := fw.Run()
		require.NoError(t, err)
	})
}

func TestPipe_Failures(t *testing.T) {
	t.Run("When first receive returns error, should return it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeError := errors.New(gofakeit.SentenceSimple())

		inMock := NewMockStream(ctrl)
		inMock.EXPECT().
			Receive().
			Return(nil, fakeError)

		outMock := NewMockStream(ctrl)

		fw := Pipe{
			In:  inMock,
			Out: outMock,
		}
		err := fw.Run()
		require.Equal(t, fakeError, err)
	})

	t.Run("When has a good message and then a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeError := errors.New(gofakeit.SentenceSimple())

		fakeMessage := new(any.Any)
		gofakeit.Struct(fakeMessage)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)

		inMock := NewMockStream(ctrl)
		inMock.EXPECT().
			Receive().
			Return(fakeMessage, nil)
		inMock.EXPECT().
			Receive().
			Return(nil, fakeError)

		outMock := NewMockStream(ctrl)
		outMock.EXPECT().Send(fakeMessage).Return(nil)

		fw := Pipe{
			In:  inMock,
			Out: outMock,
		}
		err := fw.Run()
		require.Equal(t, fakeError, err)
	})
}
