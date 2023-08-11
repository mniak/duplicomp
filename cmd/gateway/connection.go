package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EstablishedConnection struct {
	PrimaryClientConnection *grpc.ClientConn
	UseShadow               bool
	ShadowClientConnection  *grpc.ClientConn
}

func (p *EstablishedConnection) Handler(_ any, protoServer grpc.ServerStream) error {
	ctx := protoServer.Context()
	ctx, stop := context.WithCancel(ctx)
	defer stop()

	method, hasName := grpc.MethodFromServerStream(protoServer)
	if !hasName {
		return status.Errorf(codes.NotFound, "Method name could not be determined")
	}
	log.Printf("Handling method %s", method)
	defer log.Printf("Done handling method %s", method)

	ctx = copyHeadersFromIncomingToOutcoming(ctx, ctx)

	server := NewProtoStream(protoServer)
	serverObservable := ObservableFromStream(ctx, server)

	if p.UseShadow {
		go func() {
			shadow := Forwarder{
				LogName:           "Shadow",
				Method:            method,
				Server:            server,
				ServerObservable:  serverObservable,
				InboundConnection: p.ShadowClientConnection,
				DiscardResponses:  true,
			}
			err := shadow.Run(ctx)
			if err != nil {
				log.Printf("[Shadow] Failure: %s", err.Error())
			}
		}()
	}

	primary := Forwarder{
		Method:            method,
		Server:            server,
		ServerObservable:  serverObservable,
		InboundConnection: p.PrimaryClientConnection,
	}
	err := primary.Run(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}
