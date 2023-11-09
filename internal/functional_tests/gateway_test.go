package functional_tests

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp"
	"github.com/mniak/duplicomp/internal/samples"
	"github.com/mniak/duplicomp/internal/samples/grpc"
	"github.com/mniak/duplicomp/log2"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGateway_HappyPath(t *testing.T) {
	PRIMARY_PORT := gofakeit.IntRange(63000, 65000)
	SHADOW_PORT := PRIMARY_PORT + 1
	GATEWAY_PORT := PRIMARY_PORT + 2

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rootLogger := log2.Sub(log2.FromWriter(os.Stdout), "    ")
	mainLogger := log2.Sub(rootLogger, "TEST ")

	fakePingMessage := gofakeit.SentenceSimple()
	fakeMethod := gofakeit.BuzzWord()

	// ------- Primary server --------
	var fakePrimaryPong grpc.Pong
	gofakeit.Struct(&fakePrimaryPong)
	mockPrimaryHandler := samples.NewMockServerHandler(ctrl)
	mockPrimaryHandler.EXPECT().
		Handle(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
			require.Equal(t, fakePingMessage, *ping.Message)
			return &fakePrimaryPong, nil
		})
	primary, err := samples.RunServer(
		samples.WithLogger(log2.Sub(rootLogger, "PRIM ")),
		samples.WithPort(PRIMARY_PORT),
		samples.WithHandler(mockPrimaryHandler),
	)
	require.NoError(t, err)
	defer primary.Stop()

	// ------- Shadow server --------
	var fakeShadowPong grpc.Pong
	gofakeit.Struct(&fakeShadowPong)
	mockShadowHandler := samples.NewMockServerHandler(ctrl)
	mockShadowHandler.EXPECT().
		Handle(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
			require.Equal(t, fakePingMessage, *ping.Message)
			return &fakeShadowPong, nil
		})
	secondary, err := samples.RunServer(
		samples.WithLogger(log2.Sub(rootLogger, "SHAD ")),
		samples.WithPort(SHADOW_PORT),
		samples.WithHandler(mockShadowHandler),
	)
	require.NoError(t, err)
	defer secondary.Stop()

	// ------- Gateway --------
	time.Sleep(1 * time.Second)
	mockComparator := NewMockComparator(ctrl)
	mockComparator.EXPECT().Compare(fakeMethod, gomock.Any(), nil, gomock.Any(), nil)
	mockComparator.EXPECT().Compare(fakeMethod, gomock.Any(), io.EOF, gomock.Any(), io.EOF)

	mainLogger.Print("Starting gateway")
	stopGw, err := duplicomp.StartNewGateway(
		fmt.Sprintf(":%d", GATEWAY_PORT),
		duplicomp.Target{
			Address: fmt.Sprintf(":%d", PRIMARY_PORT),
		},
		duplicomp.Target{
			Address: fmt.Sprintf(":%d", SHADOW_PORT),
		},
		mockComparator,
	)
	defer stopGw.GracefulStop()
	require.NoError(t, err)

	// ------- Client --------
	time.Sleep(1 * time.Second)
	pong, err := samples.RunSendPing(
		fakePingMessage,
		samples.WithLogger(log2.Sub(rootLogger, "CLIE ")),
		samples.WithPort(GATEWAY_PORT),
	)

	time.Sleep(1 * time.Second)
	require.NoError(t, err)
	require.Equal(t, fakePrimaryPong.Reply, pong.Reply)
	require.Equal(t, fakePrimaryPong.ServedBy, pong.ServedBy)
	// -------- THE END ----------
}

func ptr[T any](t T) *T {
	return &t
}
