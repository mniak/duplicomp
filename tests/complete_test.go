package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/gateway"
	"github.com/mniak/duplicomp/internal/samples"
	"github.com/mniak/duplicomp/internal/samples/grpc"
	"github.com/mniak/duplicomp/log2"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestComplete(t *testing.T) {
	const PRIMARY_PORT = 9999
	const SHADOW_PORT = 8888
	const GATEWAY_PORT = 9000

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	rootLogger := log2.Sub(log2.FromWriter(os.Stdout), "    ")
	mainLogger := log2.Sub(rootLogger, "TEST ")

	fakePingMessage := gofakeit.SentenceSimple()

	// ------- Primary server --------

	var fakePrimaryPong grpc.Pong
	gofakeit.Struct(&fakePrimaryPong)
	mockPrimaryHandler := NewMockServerHandler(ctrl)
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
	mockShadowHandler := NewMockServerHandler(ctrl)
	mockShadowHandler.EXPECT().
		Handle(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, ping *grpc.Ping) (*grpc.Pong, error) {
			require.Equal(t, fakePingMessage, *ping.Message)
			return &fakeShadowPong, nil
		}).Times(0)
	secondary, err := samples.RunServer(
		samples.WithLogger(log2.Sub(rootLogger, "SHAD ")),
		samples.WithPort(SHADOW_PORT),
		samples.WithHandler(mockShadowHandler),
	)
	require.NoError(t, err)
	defer secondary.Stop()

	// Gateway
	time.Sleep(1 * time.Second)
	go func() {
		mainLogger.Println("Starting gateway")
		gateway.RunGateway(ctx, gateway.GatewayParams{
			ListenPort:    GATEWAY_PORT,
			PrimaryTarget: fmt.Sprintf(":%d", PRIMARY_PORT),
			ShadowTarget:  fmt.Sprintf(":%d", SHADOW_PORT),
		})
	}()

	// Client
	time.Sleep(1 * time.Second)
	pong, err := samples.RunSendPing(
		fakePingMessage,
		samples.WithLogger(log2.Sub(rootLogger, "CLIE ")),
		samples.WithPort(GATEWAY_PORT),
	)
	require.NoError(t, err)
	require.Equal(t, fakePrimaryPong.Reply, pong.Reply)
	require.Equal(t, fakePrimaryPong.ServedBy, pong.ServedBy)
}

func ptr[T any](t T) *T {
	return &t
}
