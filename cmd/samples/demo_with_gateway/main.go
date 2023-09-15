package main

import (
	"log"
	"os"
	"time"

	"github.com/mniak/duplicomp/internal/samples"
	"github.com/samber/lo"
)

func main() {
	// ctx, stop := context.WithCancel(context.Background())
	// defer stop()

	logger := log.New(os.Stdout, "[Main] ", 0)

	const PRIMARY_PORT = 9999
	const SHADOW_PORT = 8888
	const GATEWAY_PORT = 9000

	// Primary server
	logger.Println("Starting primary server")
	primary := lo.Must(samples.RunServer(samples.WithPort(PRIMARY_PORT)))
	time.Sleep(3 * time.Second)
	defer primary.Stop()

	// // Shadow server
	// logger.Println("Starting shadow server")
	// secondary := lo.Must(samples.RunServer(samples.WithPort(SHADOW_PORT)))
	// time.Sleep(3 * time.Second)
	// defer secondary.Stop()

	// // Gateway
	// go func() {
	// 	logger.Println("Starting gateway")
	// 	gateway.RunGateway(ctx, gateway.ProxyParams{
	// 		ListenPort:    GATEWAY_PORT,
	// 		PrimaryTarget: fmt.Sprintf(":%d", PRIMARY_PORT),
	// 		ShadowTarget:  fmt.Sprintf(":%d", SHADOW_PORT),
	// 	})
	// }()
	// time.Sleep(3 * time.Second)

	// // Client
	logger.Println("Sending PING")
	lo.Must0(samples.RunSendPing(samples.WithPort(PRIMARY_PORT)))
}
