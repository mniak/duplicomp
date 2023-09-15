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
	"github.com/mniak/duplicomp/log2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestComplete(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	const PRIMARY_PORT = 9999
	const SHADOW_PORT = 8888
	const GATEWAY_PORT = 9000

	rootLogger := log2.Sub(log2.FromWriter(os.Stdout), "    ")
	mainLogger := log2.Sub(rootLogger, "TEST ")

	// Primary server
	primary := lo.Must(samples.RunServer(
		samples.WithLogger(log2.Sub(rootLogger, "PRIM ")),
		samples.WithPort(PRIMARY_PORT),
	))
	defer primary.Stop()

	// Shadow server
	secondary := lo.Must(samples.RunServer(
		samples.WithLogger(log2.Sub(rootLogger, "SHAD ")),
		samples.WithPort(SHADOW_PORT),
	))
	defer secondary.Stop()

	// Gateway
	time.Sleep(1 * time.Second)
	go func() {
		mainLogger.Println("Starting gateway")
		gateway.RunGateway(ctx, gateway.ProxyParams{
			ListenPort:    GATEWAY_PORT,
			PrimaryTarget: fmt.Sprintf(":%d", PRIMARY_PORT),
			ShadowTarget:  fmt.Sprintf(":%d", SHADOW_PORT),
		})
	}()

	// Client
	time.Sleep(1 * time.Second)
	message := gofakeit.Sentence(8)
	pong := lo.Must(samples.RunSendPing(
		message,
		samples.WithLogger(log2.Sub(rootLogger, "CLIE ")),
		samples.WithPort(PRIMARY_PORT),
	))

	require.Equal(t, message, *pong.Reply)
}
