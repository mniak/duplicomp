package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func StartProxy(params ProxyParams) (*Proxy, error) {
	proxy := Proxy{
		UseShadow: params.ShadowTarget != "",
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", params.ListenPort))
	if err != nil {
		return nil, err
	}
	defer lis.Close()

	clientCredentials := insecure.NewCredentials()
	primaryClientConn, err := grpc.Dial(params.PrimaryTarget,
		grpc.WithTransportCredentials(clientCredentials),
		grpc.WithUserAgent("duplicomp-gateway/primary/0.0.1"),
	)
	if err != nil {
		return nil, err
	}
	defer primaryClientConn.Close()
	proxy.PrimaryClientConnection = primaryClientConn

	if proxy.UseShadow {
		shadowClientConn, err := grpc.Dial(params.ShadowTarget,
			grpc.WithTransportCredentials(clientCredentials),
			grpc.WithUserAgent("duplicomp-gateway/shadow/0.0.1"),
		)
		if err != nil {
			return nil, err
		}
		defer primaryClientConn.Close()
		proxy.ShadowClientConnection = shadowClientConn
	}

	server := grpc.NewServer(grpc.UnknownServiceHandler(proxy.Handler))
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "DummyService",
		HandlerType: (*any)(nil),
	}, nil)

	err = server.Serve(lis)
	if err != nil {
		return nil, err
	}

	return &proxy, nil
}
