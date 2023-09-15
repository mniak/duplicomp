package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mniak/duplicomp/internal/gateway"
	"github.com/mniak/duplicomp/internal/samples"
	"github.com/samber/lo"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	const PRIMARY_PORT = 9999
	const SHADOW_PORT = 8888
	const GATEWAY_PORT = 9000

	// Primary server
	log.Println("Starting primary server")
	primary := lo.Must(samples.RunServer(samples.WithPort(PRIMARY_PORT)))
	defer primary.Stop()

	// Shadow server

	log.Println("Starting shadow server")
	secondary := lo.Must(samples.RunServer(samples.WithPort(SHADOW_PORT)))
	defer secondary.Stop()

	// Gateway
	time.Sleep(1 * time.Second)
	go func() {
		log.Println("Starting gateway")
		gateway.RunGateway(ctx, gateway.ProxyParams{
			ListenPort:    GATEWAY_PORT,
			PrimaryTarget: fmt.Sprintf(":%d", PRIMARY_PORT),
			ShadowTarget:  fmt.Sprintf(":%d", SHADOW_PORT),
		})
	}()

	// Client
	time.Sleep(1 * time.Second)
	log.Println("Sending PING")
	lo.Must0(samples.RunSendPing(samples.WithPort(GATEWAY_PORT)))
}
