package duplicomp

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	pbany "github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

var _ Stream = &StreamWithShadow{}

func TestStreamWithShadow_Send(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeMessage := new(pbany.Any)
		gofakeit.Struct(fakeMessage)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)

		mockStream := NewMockStream(ctrl)
		mockShadow := NewMockStream(ctrl)
		mockLogger := NewMockShadowLogger(ctrl)

		var wg sync.WaitGroup
		mockStream.EXPECT().Send(fakeMessage).Return(nil)
		wg.Add(1)
		mockShadow.EXPECT().Send(fakeMessage).Do(func(_ any) {
			time.Sleep(50 * time.Millisecond)
			wg.Done()
		}).Return(nil)

		sut := StreamWithShadow{
			inner:  mockStream,
			shadow: mockShadow,
			logger: mockLogger,
		}

		startTime := time.Now()
		err := sut.Send(fakeMessage)
		elapsed1 := time.Now().Sub(startTime)

		wg.Wait()
		elapsed2 := time.Now().Sub(startTime)

		require.NoError(t, err)
		assert.InDelta(t, 0, elapsed1.Milliseconds(), 5)
		assert.InDelta(t, 50, elapsed2.Milliseconds(), 5)
	})
	t.Run("When primary fails, should not call shadow", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeError := errors.New(gofakeit.SentenceSimple())

		fakeMessage := new(pbany.Any)
		gofakeit.Struct(fakeMessage)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)

		mockStream := NewMockStream(ctrl)
		mockShadow := NewMockStream(ctrl)
		mockLogger := NewMockShadowLogger(ctrl)

		mockStream.EXPECT().Send(fakeMessage).Return(fakeError)

		sut := StreamWithShadow{
			inner:  mockStream,
			shadow: mockShadow,
			logger: mockLogger,
		}

		err := sut.Send(fakeMessage)

		require.ErrorIs(t, err, fakeError)
	})
	t.Run("When shadow fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeMessage := new(pbany.Any)
		gofakeit.Struct(fakeMessage)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)

		mockStream := NewMockStream(ctrl)
		mockShadow := NewMockStream(ctrl)
		mockLogger := NewMockShadowLogger(ctrl)

		var wg sync.WaitGroup
		mockStream.EXPECT().Send(fakeMessage).Return(nil)
		wg.Add(1)
		mockShadow.EXPECT().Send(fakeMessage).Do(func(_ any) {
			time.Sleep(50 * time.Millisecond)
			wg.Done()
		}).Return(nil)

		sut := StreamWithShadow{
			inner:  mockStream,
			shadow: mockShadow,
			logger: mockLogger,
		}

		startTime := time.Now()
		err := sut.Send(fakeMessage)
		elapsed1 := time.Now().Sub(startTime)

		wg.Wait()
		elapsed2 := time.Now().Sub(startTime)

		require.NoError(t, err)
		assert.InDelta(t, 0, elapsed1.Milliseconds(), 5)
		assert.InDelta(t, 50, elapsed2.Milliseconds(), 5)
	})
}
