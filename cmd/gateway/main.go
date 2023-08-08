package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/samber/lo"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	var listenPort int
	flag.IntVar(&listenPort, "listen-port", 9090, "TCP port to listen")
	var connectPort int
	flag.IntVar(&connectPort, "connect-port", 9000, "TCP port to connect")
	flag.Parse()

	lis := lo.Must(net.Listen("tcp", fmt.Sprintf(":%d", listenPort)))
	defer lis.Close()

	cliconn := lo.Must(grpc.Dial(fmt.Sprintf(":%d", connectPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent("duplicomp-gateway/0.0.1"),
	))
	defer lis.Close()

	proxy := Proxy{
		ClientConnection: cliconn,
	}

	server := grpc.NewServer(grpc.UnknownServiceHandler(proxy.Handler))
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "DummyService",
		HandlerType: (*any)(nil),
	}, nil)
	lo.Must0(server.Serve(lis))
}

type Proxy struct {
	ClientConnection *grpc.ClientConn
}

func (p *Proxy) Handler(srv interface{}, serverStream grpc.ServerStream) error {
	ctx, stop := context.WithCancel(serverStream.Context())
	defer stop()

	method, hasName := grpc.MethodFromServerStream(serverStream)
	if !hasName {
		return status.Errorf(codes.NotFound, "Method name could not be determined")
	}

	ctx = copyHeadersFromIncomingToOutcoming(ctx)
	log.Printf("Handling method %s", method)
	defer log.Printf("Done handling method %s", method)

	clientStream, err := p.ClientConnection.NewStream(ctx, &grpc.StreamDesc{}, method)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer clientStream.CloseSend()

	var wg sync.WaitGroup
	wg.Add(2)
	var merr error
	go func() {
		err := p.clientToServer(ctx, clientStream, serverStream)
		if err != nil {
			log.Printf("Client-to-Server failed: %s", err.Error())
			multierr.AppendInto(&merr, err)
		}
		wg.Done()
	}()
	go func() {
		err := p.serverToClient(ctx, serverStream, clientStream)
		if err != nil {
			log.Printf("Server-to-Client failed: %s", err.Error())
			multierr.AppendInto(&merr, err)
		}
		wg.Done()
	}()
	wg.Wait()
	if merr != nil {
		return status.Errorf(codes.Internal, merr.Error())
	}

	return nil
	// return status.Errorf(codes.Unavailable, "Proxy failed")
}

func (p *Proxy) serverToClient(ctx context.Context, srv grpc.ServerStream, cli grpc.ClientStream) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg := new(emptypb.Empty)
			err := srv.RecvMsg(msg)

			if errors.Is(err, io.EOF) {
				return nil
			} else if err != nil {
				return err
			}

			err = cli.SendMsg(msg)
			if err != nil {
				return err
			}
		}
	}
}

func (p *Proxy) clientToServer(ctx context.Context, cli grpc.ClientStream, srv grpc.ServerStream) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg := new(emptypb.Empty)
			err := cli.RecvMsg(msg)

			if errors.Is(err, io.EOF) {
				return nil
			} else if err != nil {
				return err
			}

			err = srv.SendMsg(msg)
			if err != nil {
				return err
			}
		}
	}
}

func copyHeadersFromIncomingToOutcoming(ctx context.Context) context.Context {
	meta, hasMeta := metadata.FromIncomingContext(ctx)
	if hasMeta {
		for k, vals := range meta {
			for _, v := range vals {
				ctx = metadata.AppendToOutgoingContext(ctx, k, v)
			}
		}
	}
	return ctx
}
