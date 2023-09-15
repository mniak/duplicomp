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
	"google.golang.org/protobuf/proto"
)

type Forwarder struct {
	LogName    string
	Method     string
	Downstream Stream
	Upstream   Stream
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
	protoIn, err := f.Upstream.NewStream(ctx, &grpc.StreamDesc{}, f.Method)
	if err != nil {
		return err
	}
	in := InOutStream(StreamFromProtobuf(protoIn))
	clientObs := ObservableFromStream(ctx, in)

	// Receive from server and forward to client
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := f.forwardMessages(f.DownstreamObservable, in)
		multierr.AppendInto(&combinedErrors, err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer protoIn.CloseSend()
		err := f.forwardMessages(clientObs, f.Downstream)
		multierr.AppendInto(&combinedErrors, err)
	}()

	wg.Wait()
	return combinedErrors
}

func (f *Forwarder) forwardMessages(from rxgo.Observable, to Stream) error {
	var err error
	for item := range from.Observe() {
		if item.Error() {
			return item.E
		}
		msg := item.V.(proto.Message)
		err = to.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func ObservableFromProtoStream(ctx context.Context, protoStream iProtoStream) rxgo.Observable {
	stream := InOutStream(StreamFromProtobuf(protoStream))
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
