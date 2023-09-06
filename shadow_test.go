package duplicomp

import (
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	pbany "github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

var _ Stream = &ShadowStream{}

func TestShadowLogger_Send(t *testing.T) {
	t.Run("Should send both concurrently and return even after secondary sends", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeMessage := new(pbany.Any)
		gofakeit.Struct(fakeMessage)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)

		mockPrimaryStream := NewMockStream(ctrl)
		mockShadowStream := NewMockStream(ctrl)
		mockShadowLogger := NewMockShadowLogger(ctrl)

		var wg sync.WaitGroup
		mockPrimaryStream.EXPECT().Send(fakeMessage).Return(nil)
		wg.Add(1)
		mockShadowStream.EXPECT().Send(fakeMessage).Do(func(_ any) {
			time.Sleep(50 * time.Millisecond)
			wg.Done()
		}).Return(nil)

		sut := ShadowStream{
			primaryStream: mockPrimaryStream,
			shadowStream:  mockShadowStream,
			shadowLogger:  mockShadowLogger,
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

	// t.Run("When shadow fails", func(t *testing.T) {
	// 	ctrl := gomock.NewController(t)
	// 	defer ctrl.Finish()

	// 	mockPrimaryStream := NewMockStream(ctrl)
	// 	mockShadowStream := NewMockStream(ctrl)
	// 	mockShadowLogger := NewMockShadowLogger(ctrl)

	// 	fakePrimaryError := errors.New(gofakeit.SentenceSimple())
	// 	fakeShadowError := errors.New(gofakeit.SentenceSimple())

	// 	mockShadowLogger.recorder.LogSendError(fakeShadowError)

	// 	fakeMessage := new(pbany.Any)
	// 	gofakeit.Struct(fakeMessage)
	// 	require.NotNil(t, fakeMessage)
	// 	require.NotEmpty(t, fakeMessage)

	// 	sut := ShadowStream{
	// 		primaryStream: mockPrimaryStream,
	// 		shadowStream:  mockShadowStream,
	// 		shadowLogger:  mockShadowLogger,
	// 	}

	// 	resultError := sut.Send(fakeMessage)
	// 	assert.Equal(t, fakePrimaryError, resultError)
	// })
}
