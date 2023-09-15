package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/gateway"
	"github.com/mniak/duplicomp/internal/samples"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestComplete(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	logger := log.New(os.Stdout, "[Main] ", 0)

	const PRIMARY_PORT = 9999
	const SHADOW_PORT = 8888
	const GATEWAY_PORT = 9000

	// Primary server
	logger.Println("Starting primary server")
	primary := lo.Must(samples.RunServer(
		samples.WithName("Primary"),
		samples.WithPort(PRIMARY_PORT),
	))
	defer primary.Stop()

	// Shadow server
	logger.Println("Starting shadow server")
	secondary := lo.Must(samples.RunServer(
		samples.WithName("Shadow"),
		samples.WithPort(SHADOW_PORT),
	))
	defer secondary.Stop()

	// Gateway
	time.Sleep(3 * time.Second)
	go func() {
		logger.Println("Starting gateway")
		gateway.RunGateway(ctx, gateway.ProxyParams{
			ListenPort:    GATEWAY_PORT,
			PrimaryTarget: fmt.Sprintf(":%d", PRIMARY_PORT),
			ShadowTarget:  fmt.Sprintf(":%d", SHADOW_PORT),
		})
	}()

	// Client
	time.Sleep(3 * time.Second)
	logger.Println("Sending PING")
	message := gofakeit.Sentence(8)
	pong := lo.Must(samples.RunSendPing(
		message,
		samples.WithName("Client"),
		samples.WithPort(PRIMARY_PORT),
	))

	require.Equal(t, message, *pong.OriginalMessage)
}
