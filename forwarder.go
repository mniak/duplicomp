package duplicomp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Forwarder struct {
	LogName           string
	Method            string
	Server            Stream
	ServerObservable  rxgo.Observable
	InboundConnection *grpc.ClientConn
	DiscardResponses  bool
}

func (f *Forwarder) Run(ctx context.Context) (err error) {
	defer func() {
		data := recover()
		if data != nil {
			err = fmt.Errorf("forwarder panicked: %+v", data)
		}
	}()

	var wg sync.WaitGroup
	var combinedErrors error
	protoIn, err := f.InboundConnection.NewStream(ctx, &grpc.StreamDesc{}, f.Method)
	if err != nil {
		return err
	}
	in := NewProtoStream(protoIn)
	clientObs := ObservableFromStream(ctx, in)

	// Receive from server and forward to client
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := f.forwardMessages(f.ServerObservable, in)
		multierr.AppendInto(&combinedErrors, err)
	}()

	// Receive from client and forward to server
	if !f.DiscardResponses {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer protoIn.CloseSend()
			err := f.forwardMessages(clientObs, f.Server)
			multierr.AppendInto(&combinedErrors, err)
		}()
	}

	wg.Wait()
	return combinedErrors
}

func (f *Forwarder) forwardMessages(from rxgo.Observable, to Stream) error {
	var err error
	for item := range from.Observe() {
		if item.Error() {
			err = item.E
			break
		}
		msg := item.V.(Message)
		err = to.Send(msg)
		if err != nil {
			break
		}
	}
	return err
}

func ObservableFromProtoStream(ctx context.Context, protoStream iProtoStream) rxgo.Observable {
	stream := NewProtoStream(protoStream)
	obs := ObservableFromStream(ctx, stream)
	return obs
}

func ObservableFromStream(ctx context.Context, stream Stream) rxgo.Observable {
	defer func() {
		data := recover()
		if data != nil {
			log.Printf("[ReceiveFromStream] PANIC: %+v", data)
		}
	}()

	obs := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		for {
			select {
			case <-ctx.Done():
				next <- rxgo.Error(ctx.Err())
				return
			default:
				msg, err := stream.Receive()

				if errors.Is(err, io.EOF) {
					return
				}
				if err != nil {
					next <- rxgo.Error(err)
					return
				}

				next <- rxgo.Of(msg)
			}
		}
	}})
	return obs
}

func copyHeadersFromIncomingToOutcoming(in, out context.Context) context.Context {
	meta, hasMeta := metadata.FromIncomingContext(in)
	if hasMeta {
		for k, vals := range meta {
			for _, v := range vals {
				in = metadata.AppendToOutgoingContext(out, k, v)
			}
		}
	}
	return in
}

type Forwarder2 struct {
	InboundStream Stream
}

func (fw *Forwarder2) Run() error {
	fw.InboundStream.Receive()
	fw.InboundStream.Receive()
	return nil
}
