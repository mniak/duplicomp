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

	"github.com/reactivex/rxgo/v2"
	"github.com/samber/lo"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	var listenPort int
	flag.IntVar(&listenPort, "listen-port", 9091, "TCP port to listen")
	var primaryTarget string
	flag.StringVar(&primaryTarget, "target", ":9001", "Connection target")
	var shadowTarget string
	flag.StringVar(&shadowTarget, "shadow-target", "", "Shadow connection target")
	flag.Parse()

	useShadow := shadowTarget != ""

	lis := lo.Must(net.Listen("tcp", fmt.Sprintf(":%d", listenPort)))
	defer lis.Close()

	clientCredentials := insecure.NewCredentials()

	primaryClientConn := lo.Must(grpc.Dial(primaryTarget,
		grpc.WithTransportCredentials(clientCredentials),
		grpc.WithUserAgent("duplicomp-gateway/primary/0.0.1"),
	))
	defer primaryClientConn.Close()

	proxy := Proxy{
		PrimaryClientConnection: primaryClientConn,
	}
	if useShadow {
		shadowClientConn := lo.Must(grpc.Dial(shadowTarget,
			grpc.WithTransportCredentials(clientCredentials),
			grpc.WithUserAgent("duplicomp-gateway/shadow/0.0.1"),
		))
		defer primaryClientConn.Close()
		proxy.ShadowClientConnection = shadowClientConn
	}

	server := grpc.NewServer(grpc.UnknownServiceHandler(proxy.Handler))
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "DummyService",
		HandlerType: (*any)(nil),
	}, nil)
	lo.Must0(server.Serve(lis))
}

type Proxy struct {
	PrimaryClientConnection *grpc.ClientConn
	ShadowClientConnection  *grpc.ClientConn
}

func (p *Proxy) Handler(_ any, server grpc.ServerStream) error {
	ctx, stop := context.WithCancel(server.Context())
	defer stop()

	method, hasName := grpc.MethodFromServerStream(server)
	if !hasName {
		return status.Errorf(codes.NotFound, "Method name could not be determined")
	}

	ctx = copyHeadersFromIncomingToOutcoming(ctx)
	log.Printf("Handling method %s", method)
	defer log.Printf("Done handling method %s", method)

	primaryClient, err := p.PrimaryClientConnection.NewStream(ctx, &grpc.StreamDesc{}, method)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer primaryClient.CloseSend()

	osrv := p.receiveFromServer(ctx, server)
	// var shadowClient grpc.ClientStream
	// if p.ShadowClientConnection != nil {
	// 	shadowClient, err = p.ShadowClientConnection.NewStream(ctx, &grpc.StreamDesc{}, method)
	// 	if err != nil {
	// 		return status.Error(codes.Internal, err.Error())
	// 	}
	// 	defer shadowClient.CloseSend()
	// }

	var wg sync.WaitGroup
	wg.Add(2)
	var merr error
	go func() {
		defer wg.Done()
		err := p.clientToServer(ctx, primaryClient, server)
		if err != nil {
			log.Printf("Client-to-Server failed: %s", err.Error())
			multierr.AppendInto(&merr, err)
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		for item := range osrv.Observe() {
			if item.Error() {
				err = item.E
				break
			}
			msg := item.V.(proto.Message)

			err = primaryClient.SendMsg(msg)
			if err != nil {
				break
			}
		}

		if err != nil {
			log.Printf("Server-to-Client failed: %s", err.Error())
			multierr.AppendInto(&merr, err)
		}
	}()
	wg.Wait()
	if merr != nil {
		return status.Errorf(codes.Internal, merr.Error())
	}

	return nil
	// return status.Errorf(codes.Unavailable, "Proxy failed")
}

func (p *Proxy) receiveFromServer(ctx context.Context, srv grpc.ServerStream) rxgo.Observable {
	obs := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		for {
			select {
			case <-ctx.Done():
				next <- rxgo.Error(ctx.Err())
				return
			default:
				msg := new(emptypb.Empty)
				err := srv.RecvMsg(msg)

				if errors.Is(err, io.EOF) {
					return
				}
				if err != nil {
					next <- rxgo.Error(err)
					return
				}

				next <- rxgo.Of(msg)

				// err = primaryCli.SendMsg(msg)
				// if err != nil {
				// 	return err
				// }

				// // TODO: do something to avoid that a failure sending to the shadow to affect the connection
				// if shadowCli != nil {
				// 	err = shadowCli.SendMsg(msg)
				// 	if err != nil {
				// 		return err
				// 	}
				// }
			}
		}
	}})
	return obs
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
