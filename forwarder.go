package ps121

import (
	"context"
	"errors"
	"io"
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
		obs := ObservableFromStream(ctx, f.Downstream)
		err := f.forwardMessages(obs, f.Upstream)
		errorChan <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		obs := ObservableFromStream(ctx, f.Upstream)
		err := f.forwardMessages(obs, f.Downstream)
		errorChan <- err
	}()

	wg.Wait()
	return combinedErrors
}

func (f *Forwarder) forwardMessages(from rxgo.Observable, to OutputStream) error {
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

func ObservableFromStream(ctx context.Context, stream InputStream) rxgo.Observable {
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
