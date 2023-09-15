package main

import (
	"context"
	"log"

	"github.com/mniak/duplicomp/internal/gateway"
	"github.com/mniak/duplicomp/internal/samples"
	"github.com/samber/lo"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	// Primary server
	log.Println("Starting primary server")
	primary := lo.Must(samples.RunServer(samples.WithPort(9999)))
	defer primary.Stop()

	// Shadow server
	log.Println("Starting shadow server")
	secondary := lo.Must(samples.RunServer(samples.WithPort(8888)))
	defer secondary.Stop()

	// Gateway
	go func() {
		log.Println("Starting gateway")
		gateway.RunGateway(ctx, gateway.ProxyParams{
			ListenPort:    9000,
			PrimaryTarget: ":9999",
			ShadowTarget:  ":8888",
		})
	}()

	// Client
	log.Println("Sending PING")
	lo.Must0(samples.RunSendPing(samples.WithPort(9000)))
}
