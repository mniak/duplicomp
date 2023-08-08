package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	var listenPort int
	flag.IntVar(&listenPort, "listen-port", 8089, "TCP port to listen")
	var connectPort int
	flag.IntVar(&connectPort, "connect-port", 8099, "TCP port to connect")
	flag.Parse()

	lis := lo.Must(net.Listen("tcp", fmt.Sprintf("localhost:%d", listenPort)))
	defer lis.Close()

	server := grpc.NewServer(grpc.UnknownServiceHandler(UnknownService))
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "MyService",
		HandlerType: (*any)(nil),
	}, nil)
	lo.Must0(server.Serve(lis))
}

func UnknownService(srv interface{}, stream grpc.ServerStream) error {
	grpc.MethodFromServerStream(stream)
	return status.Errorf(codes.Unavailable, "Proxy failed")
}
