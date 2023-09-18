package duplicomp

import (
	"context"
	"net"

	"github.com/mniak/duplicomp/internal/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayParams struct {
	ListenPort    int
	PrimaryTarget string
	ShadowTarget  string
	Comparator    Comparator
}

type LambdaConnectionHandler func(ctx context.Context, method string, serverStream Stream) error

func (lch LambdaConnectionHandler) HandleConnection(ctx context.Context, method string, serverStream Stream) error {
	return lch(ctx, method, serverStream)
}

type Gateway interface {
	Start(ctx context.Context) error
	Stop()
	Wait() error
}

type _Gateway struct {
	Listener          net.Listener
	PrimaryConnection TargetConnection
	ShadowConnection  TargetConnection
	Comparator        Comparator

	grpcServer GRPCServer
}

func NewGateway(listenAddr, primaryTarget, shadowTarget string) (Gateway, error) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	primaryConnection, err := ConnectionSpec{Address: primaryTarget}.Connect()
	if err != nil {
		return nil, err
	}

	shadowConnection, err := ConnectionSpec{Address: shadowTarget}.Connect()
	if err != nil {
		return nil, err
	}

	// var comparator LogComparator

	gw := _Gateway{
		PrimaryConnection: primaryConnection,
		ShadowConnection:  shadowConnection,
		Listener:          listener,
		// Comparator:        comparator,
	}
	return &gw, nil
}

func (gw *_Gateway) Start(ctx context.Context) error {
	gw.grpcServer = GRPCServer{
		ConnectionHandler: LambdaConnectionHandler(func(ctx context.Context, method string, serverStream Stream) error {
			primaryUpstream, err := gw.PrimaryConnection.Stream(ctx, method)
			if err != nil {
				return err
			}

			shadowUpstream, err := gw.ShadowConnection.Stream(ctx, method)
			if err != nil {
				return err
			}

			dualStream := StreamWithShadow{
				Primary: primaryUpstream,
				Shadow:  shadowUpstream,
				Logger:  noop.Logger(),
			}

			fwd := Forwarder{
				Downstream: serverStream,
				Upstream:   &dualStream,
			}
			return fwd.Run(ctx)
		}),
	}

	err := gw.grpcServer.StartWithListener(gw.Listener)
	return err
}

func (gw *_Gateway) Stop() {
	gw.grpcServer.Stop()
}

func (gw *_Gateway) Wait() error {
	return gw.grpcServer.Wait()
}

// func establishConnections(primaryTarget, shadowTarget string) TargetConnection {
// 	specs := []ConnectionSpec{
// 		{Address: primaryTarget},
// 		{Address: shadowTarget},
// 	}

// allConnected := make([]TargetConnection, 0, len(specs)

// 	for _, spec := range specs {
// 		conn, err := spec.Connect()
// 		if err == nil {
// 			return conn
// 		}
// 	}
// }

type grpcConnection struct {
	conn *grpc.ClientConn
}

func (self *grpcConnection) Stream(ctx context.Context, method string) (Stream, error) {
	protoStream, err := self.conn.NewStream(ctx, &grpc.StreamDesc{}, method)
	if err != nil {
		return nil, err
	}
	stream := InOutStream(StreamFromProtobuf(protoStream))
	return stream, nil
}

func (self *grpcConnection) Close() {
	if self.conn != nil {
		self.conn.Close()
	}
}

type ConnectionSpec struct {
	Address string
}

func (self ConnectionSpec) Connect() (TargetConnection, error) {
	clientCredentials := insecure.NewCredentials()
	conn, err := grpc.Dial(self.Address,
		grpc.WithTransportCredentials(clientCredentials),
		grpc.WithUserAgent("duplicomp-gateway/0.0.1"),
	)
	if err != nil {
		return nil, err
	}
	return &grpcConnection{
		conn: conn,
	}, nil
}
