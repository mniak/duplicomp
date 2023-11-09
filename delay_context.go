package duplicomp

import (
	"context"
	"sync"
	"time"
)

func ContextWithDelay(ctx context.Context, delay time.Duration) context.Context {
	return &_DelayContext{
		Context: ctx,
		Delay:   delay,
	}
}

type _DelayContext struct {
	Context context.Context
	Delay   time.Duration

	done         chan struct{}
	onceInitDone sync.Once
}

func (self *_DelayContext) Deadline() (deadline time.Time, ok bool) {
	deadline, hasDeadline := self.Context.Deadline()
	if !hasDeadline {
		return time.Time{}, false
	}
	return deadline.Add(self.Delay), false
}

func (self *_DelayContext) Done() <-chan struct{} {
	self.onceInitDone.Do(func() {
		self.done = make(chan struct{})
		go func() {
			<-self.Context.Done()
			time.Sleep(self.Delay)
			close(self.done)
		}()
	})

	return self.done
}

func (self *_DelayContext) Err() error {
	select {
	case <-self.Done():
		return self.Context.Err()
	default:
		return nil
	}
}

func (self *_DelayContext) Value(key any) any {
	return self.Context.Value(key)
}
