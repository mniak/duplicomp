package overflowchan

import "sync"

type Channel[T any] struct {
	dataChan     chan T
	closedSignal chan struct{}
	onceClose    sync.Once
}

const OverflowableChannelDefaultBufferSize = 2

func New[T any](bufferSize int) Channel[T] {
	if bufferSize < 1 {
		bufferSize = OverflowableChannelDefaultBufferSize
	}
	return Channel[T]{
		dataChan:     make(chan T, bufferSize),
		closedSignal: make(chan struct{}),
	}
}

func (self Channel[T]) Send(val T) {
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

func (self Channel[T]) Close() {
	self.onceClose.Do(func() {
		close(self.closedSignal)
		close(self.dataChan)
	})
}

func (self Channel[T]) Receiver() <-chan T {
	return self.dataChan
}
