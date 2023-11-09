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
		mockComparator := NewMockComparator(ctrl)

		var wg sync.WaitGroup
		mockStream.EXPECT().Send(fakeMessage).Return(nil)
		wg.Add(1)
		mockShadow.EXPECT().Send(fakeMessage).Do(func(_ any) {
			time.Sleep(50 * time.Millisecond)
			wg.Done()
		}).Return(nil)

		sut := StreamWithShadow{
			Primary:    mockStream,
			Shadow:     mockShadow,
			Comparator: mockComparator,
		}

		startTime := time.Now()
		err := sut.Send(fakeMessage)
		elapsed1 := time.Now().Sub(startTime)

		wg.Wait()
		elapsed2 := time.Now().Sub(startTime)

		require.NoError(t, err)
		assert.InDelta(t, 0, elapsed1.Milliseconds(), 5)
		assert.InDelta(t, 50, elapsed2.Milliseconds(), 5)

		// Wait for remaining calls in other goroutines
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("When primary fails, dont call shadow", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeError := errors.New(gofakeit.SentenceSimple())

		fakeMessage := new(pbany.Any)
		gofakeit.Struct(fakeMessage)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)

		mockStream := NewMockStream(ctrl)
		mockShadow := NewMockStream(ctrl)
		mockComparator := NewMockComparator(ctrl)

		mockStream.EXPECT().Send(fakeMessage).Return(fakeError)

		sut := StreamWithShadow{
			Primary:    mockStream,
			Shadow:     mockShadow,
			Comparator: mockComparator,
		}

		err := sut.Send(fakeMessage)
		require.ErrorIs(t, err, fakeError)

		// Wait for remaining calls in other goroutines
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("When shadow fails must log", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeError := errors.New(gofakeit.SentenceSimple())

		fakeMessage := new(pbany.Any)
		gofakeit.Struct(fakeMessage)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)

		mockStream := NewMockStream(ctrl)
		mockShadow := NewMockStream(ctrl)
		mockComparator := NewMockComparator(ctrl)

		var wg sync.WaitGroup
		mockStream.EXPECT().Send(fakeMessage).Return(nil)
		wg.Add(1)
		mockShadow.EXPECT().Send(fakeMessage).Do(func(_ any) {
			time.Sleep(50 * time.Millisecond)
			wg.Done()
		}).Return(fakeError)
		// mockLogger.EXPECT().LogSendFailure(fakeError)

		sut := StreamWithShadow{
			Primary:    mockStream,
			Shadow:     mockShadow,
			Comparator: mockComparator,
		}

		startTime := time.Now()
		err := sut.Send(fakeMessage)
		elapsed1 := time.Now().Sub(startTime)

		wg.Wait()
		elapsed2 := time.Now().Sub(startTime)

		require.NoError(t, err)
		assert.InDelta(t, 0, elapsed1.Milliseconds(), 5)
		assert.InDelta(t, 50, elapsed2.Milliseconds(), 5)

		// Wait for remaining calls in other goroutines
		time.Sleep(100 * time.Millisecond)
	})
}

func TestStreamWithShadow_Receive(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeMessageBytes := []byte(gofakeit.SentenceSimple())
		fakeMessage := NewFakeProtoMessage(fakeMessageBytes)
		require.NotNil(t, fakeMessage)
		require.NotEmpty(t, fakeMessage)
		// fakeError := errors.New(gofakeit.SentenceSimple())
		fakeError := errors.New("fake error primary")

		fakeShadowMessageBytes := []byte(gofakeit.SentenceSimple())
		fakeShadowMessage := NewFakeProtoMessage(fakeShadowMessageBytes)
		require.NotNil(t, fakeShadowMessage)
		require.NotEmpty(t, fakeShadowMessage)
		// fakeShadowError := errors.New(gofakeit.SentenceSimple())
		fakeShadowError := errors.New("fake error shadow")

		require.NotEqual(t, fakeMessage, fakeShadowMessage)
		require.NotEqual(t, &fakeMessage, &fakeShadowMessage)

		mockStream := NewMockStream(ctrl)
		mockShadow := NewMockStream(ctrl)
		mockComparator := NewMockComparator(ctrl)

		var wg sync.WaitGroup
		wg.Add(1)
		mockStream.EXPECT().Receive().Return(fakeMessage, fakeError)
		mockShadow.EXPECT().Receive().
			Do(func() {
				time.Sleep(50 * time.Millisecond)
				wg.Done()
			}).
			Return(fakeShadowMessage, fakeShadowError)
		mockComparator.EXPECT().Compare(fakeMessageBytes, fakeError, fakeShadowMessageBytes, fakeShadowError)

		sut := StreamWithShadow{
			Primary:    mockStream,
			Shadow:     mockShadow,
			Comparator: mockComparator,
		}

		startTime := time.Now()
		msg, receiveErr := sut.Receive()
		receiveDuration := time.Now().Sub(startTime)

		wg.Wait()
		receiveShadowDuration := time.Now().Sub(startTime)

		require.Equal(t, fakeError, receiveErr)
		assert.Equal(t, fakeMessage, msg)

		assert.InDelta(t, 0, receiveDuration.Milliseconds(), 5)
		assert.InDelta(t, 50, receiveShadowDuration.Milliseconds(), 5)

		// Wait for remaining calls in other goroutines
		time.Sleep(500 * time.Millisecond)
	})

	// 	t.Run("When primary receive fails, dont call shadow", func(t *testing.T) {
	// 		ctrl := gomock.NewController(t)
	// 		defer ctrl.Finish()

	// 		fakeMessage := new(pbany.Any)
	// 		gofakeit.Struct(fakeMessage)
	// 		require.NotNil(t, fakeMessage)
	// 		require.NotEmpty(t, fakeMessage)
	// 		fakeError := errors.New(gofakeit.SentenceSimple())

	// 		fakeShadowMessage := new(pbany.Any)
	// 		gofakeit.Struct(fakeShadowMessage)
	// 		require.NotNil(t, fakeShadowMessage)
	// 		require.NotEmpty(t, fakeShadowMessage)

	// 		require.NotEqual(t, fakeMessage, fakeShadowMessage)
	// 		require.NotEqual(t, &fakeMessage, &fakeShadowMessage)

	// 		mockStream := NewMockStream(ctrl)
	// 		mockShadow := NewMockStream(ctrl)
	// 		mockComparator := NewMockComparator(ctrl)

	// 		mockStream.EXPECT().Receive().Return(fakeMessage, fakeError)

	// 		sut := StreamWithShadow{
	// 			Primary:    mockStream,
	// 			Shadow:     mockShadow,
	// 			Comparator: mockComparator,
	// 		}

	// 		msg, err := sut.Receive()

	// 		require.ErrorIs(t, err, fakeError)
	// 		assert.Equal(t, fakeMessage, msg)

	//		// Wait for remaining calls in other goroutines
	//		time.Sleep(100 * time.Millisecond)
	//	})
}
