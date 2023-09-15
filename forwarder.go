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
	"google.golang.org/protobuf/proto"
)

type Forwarder struct {
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

	errorChan := make(chan error)
	go func() {
		for err := range errorChan {
			multierr.AppendInto(&combinedErrors, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := f.forward(ctx, f.Downstream, f.Upstream)
		errorChan <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := f.forward(ctx, f.Upstream, f.Downstream)
		errorChan <- err
	}()

	wg.Wait()
	return combinedErrors
}

func (f *Forwarder) forward(ctx context.Context, from InputStream, to OutputStream) error {
	for {
		msg, err := from.Receive()
		if err != nil {
			return err
		}
		err = to.Send(msg)
		if err != nil {
			return err
		}
	}
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
