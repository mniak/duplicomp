package main

import (
	"flag"
	"fmt"

	"github.com/mniak/duplicomp"
	"github.com/samber/lo"
)

func main() {
	var params ProxyParams
	flag.IntVar(&params.ListenPort, "listen-port", 9091, "TCP port to listen")
	flag.StringVar(&params.PrimaryTarget, "target", ":9001", "Connection target")
	flag.StringVar(&params.ShadowTarget, "shadow-target", "", "Shadow connection target")
	flag.Parse()

	lo.Must0(StartProxy(params))
}

type ProxyParams struct {
	ListenPort    int
	PrimaryTarget string
	ShadowTarget  string
}

func StartProxy(params ProxyParams) error {
	proxy := duplicomp.NewGRPCProxy(duplicomp.ProxyConfig{
		InboundConfig: duplicomp.InboundConfig{
			ListenAddress: fmt.Sprintf(":%d", params.ListenPort),
		},
		OutboundConfigs: []duplicomp.OutboundConfig{
			{
				TargetAddress: params.PrimaryTarget,
			},
			{
				TargetAddress:  params.ShadowTarget,
				IgnoreResponse: true,
			},
		},
	})

	err := proxy.Start()
	if err != nil {
		return err
	}
	err = proxy.Wait()
	if err != nil {
		return err
	}
	return nil
}
