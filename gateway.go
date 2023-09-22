package duplicomp

import (
	"context"
	"log"
	"net"

	"github.com/mniak/duplicomp/log2"
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

func (h LambdaConnectionHandler) HandleConnection(ctx context.Context, method string, serverStream Stream) error {

	return h(ctx, method, serverStream)
}

type Gateway interface {
	Start(ctx context.Context) error
	Stopper
}

type _Gateway struct {
	Listener          net.Listener
	PrimaryConnection TargetConnection
	ShadowConnection  TargetConnection
	Comparator        Comparator

	Logger log2.Logger

	grpcServer GRPCServer
}

func StartNewGateway(listenAddr, primaryTarget, shadowTarget string, cmp Comparator) (GracefulStopper, error) {
	logger := log2.Sub(log.Default(), "[Gateway] ")
	var cb PessimisticCallerback
	defer cb.Callback()

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Printf("failed to open listener at %s", listenAddr)
		return nil, err
	}
	cb.OnFailure(func() { listener.Close() })
	logger.Printf("listener open at %s", listenAddr)

	primaryConnection, err := ConnectionSpec{Address: primaryTarget}.Connect()
	if err != nil {
		logger.Printf("failed to establish primary connection to %s", primaryTarget)
		return nil, err
	}
	cb.OnFailure(func() { primaryConnection.Close() })
	logger.Printf("primary connection established to %s", primaryTarget)

	shadowConnection, err := ConnectionSpec{Address: shadowTarget}.Connect()
	if err != nil {
		logger.Printf("failed to establish shadow connection to %s", primaryTarget)
		return nil, err
	}
	cb.OnFailure(func() { shadowConnection.Close() })
	logger.Printf("shadow connection established to %s", primaryTarget)

	cb.Succeeded()

	gw := _Gateway{
		PrimaryConnection: primaryConnection,
		ShadowConnection:  shadowConnection,
		Listener:          listener,
		Logger:            logger,
		Comparator:        cmp,
	}
	gw.Start()
	return StopperFunc(gw.GracefulStop), nil
}

func (gw *_Gateway) Start() error {
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
				Primary:    primaryUpstream,
				Shadow:     shadowUpstream,
				Logger:     gw.Logger,
				Comparator: gw.Comparator,
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

func (gw *_Gateway) GracefulStop() {
	gw.grpcServer.GracefulStop()
}

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
