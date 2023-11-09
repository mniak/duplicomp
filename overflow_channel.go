package duplicomp

import "sync"

type OverflowChannel[T any] struct {
	dataChan     chan T
	closedSignal chan struct{}
	onceClose    sync.Once
}

func NewOverflowableChannel[T any](bufferSize int) OverflowChannel[T] {
	if bufferSize < 1 {
		bufferSize = OverflowableChannelDefaultBufferSize
	}
	return OverflowChannel[T]{
		dataChan:     make(chan T, bufferSize),
		closedSignal: make(chan struct{}),
	}
}

func (self OverflowChannel[T]) Send(val T) {
	select {
	case <-self.closedSignal:

	default:
		select {
		case self.dataChan <- val:
		default:
			self.Close()
		}

	}
}

func (self OverflowChannel[T]) Close() {
	self.onceClose.Do(func() {
		close(self.closedSignal)
		close(self.dataChan)
	})
}

func (self OverflowChannel[T]) Receiver() <-chan T {
	return self.dataChan
}
