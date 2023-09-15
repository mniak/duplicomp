package main

import (
	"context"
	"flag"

	"github.com/mniak/duplicomp/internal/gateway"
	"github.com/samber/lo"
)

func main() {
	var params gateway.GatewayParams
	flag.IntVar(&params.ListenPort, "listen-port", 9091, "TCP port to listen")
	flag.StringVar(&params.PrimaryTarget, "target", ":9001", "Connection target")
	flag.StringVar(&params.ShadowTarget, "shadow-target", "", "Shadow connection target")
	flag.Parse()

	lo.Must0(gateway.RunGateway(context.TODO(), params))
}
