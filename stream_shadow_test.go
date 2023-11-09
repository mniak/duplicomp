package duplicomp

import (
	"errors"
	"io"
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

func TestStreamWithShadow_Receive_Realistic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeMethod := gofakeit.BuzzWord()
	fakeMessageBytes := []byte(gofakeit.SentenceSimple())
	fakeMessage := NewFakeProtoMessage(fakeMessageBytes)
	require.NotNil(t, fakeMessage)
	require.NotEmpty(t, fakeMessage)
	fakeError := errors.New(gofakeit.SentenceSimple())

	fakeShadowMessageBytes := []byte(gofakeit.SentenceSimple())
	fakeShadowMessage := NewFakeProtoMessage(fakeShadowMessageBytes)
	require.NotNil(t, fakeShadowMessage)
	require.NotEmpty(t, fakeShadowMessage)
	fakeShadowError := errors.New(gofakeit.SentenceSimple())

	require.NotEqual(t, fakeMessage, fakeShadowMessage)
	require.NotEqual(t, &fakeMessage, &fakeShadowMessage)

	mockStream := NewMockStream(ctrl)
	mockShadow := NewMockStream(ctrl)
	mockComparator := NewMockComparator(ctrl)

	mockStream.EXPECT().
		Receive().
		Do(func() {
			time.Sleep(200 * time.Millisecond)
		}).
		Return(fakeMessage, fakeError)
	mockStream.EXPECT().Receive().Return(nil, io.EOF)

	var waitShadow sync.WaitGroup
	waitShadow.Add(2)
	mockShadow.EXPECT().
		Receive().
		Do(func() {
			time.Sleep(50 * time.Millisecond)
			waitShadow.Done()
		}).
		Return(fakeShadowMessage, fakeShadowError)
	mockShadow.EXPECT().
		Receive().
		Do(func() {
			waitShadow.Done()
		}).
		Return(nil, io.EOF)

	var waitComparator sync.WaitGroup
	waitComparator.Add(2)
	mockComparator.EXPECT().
		Compare(fakeMethod, fakeMessageBytes, fakeError, fakeShadowMessageBytes, fakeShadowError).
		Do(func(_, _, _, _, _ any) {
			waitComparator.Done()
		})
	mockComparator.EXPECT().
		Compare(fakeMethod, nil, io.EOF, nil, io.EOF).
		Do(func(_, _, _, _, _ any) {
			waitComparator.Done()
		})

	sut := StreamWithShadow{
		Method:     fakeMethod,
		Primary:    mockStream,
		Shadow:     mockShadow,
		Comparator: mockComparator,
	}

	startTime := time.Now()
	receiveMsg1, receiveErr1 := sut.Receive()
	receiveMsg2, receiveErr2 := sut.Receive()

	receivePrimaryDuration := time.Now().Sub(startTime)
	waitShadow.Wait()
	receiveShadow1Duration := time.Now().Sub(startTime)

	assert.Equal(t, fakeMessage, receiveMsg1)
	assert.Equal(t, fakeError, receiveErr1)

	assert.Nil(t, receiveMsg2)
	assert.Equal(t, io.EOF, receiveErr2)

	assert.InDelta(t, 200, receivePrimaryDuration.Milliseconds(), 5)
	assert.InDelta(t, 200+50, receiveShadow1Duration.Milliseconds(), 5)

	waitComparator.Wait()
}
