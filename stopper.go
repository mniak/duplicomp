package duplicomp

type Stopper interface {
	Stop()
}

type GracefulStopper interface {
	GracefulStop()
}

type StopperFunc func()

func (fn StopperFunc) GracefulStop() {
	if fn != nil {
		fn()
	}
}

func (fn StopperFunc) Stop() {
	fn.GracefulStop()
}

// type inlineStopWaiter struct {
// 	stop func()
// 	wait func() error
// }

// func InlineStopWaiter(stop func(), wait func() error) inlineStopWaiter {
// 	return inlineStopWaiter{
// 		stop: stop,
// 		wait: wait,
// 	}
// }

// func (sw inlineStopWaiter) Stop() {
// 	if sw.stop != nil {
// 		sw.stop()
// 	}
// }

// func (sw inlineStopWaiter) Wait() error {
// 	if sw.stop == nil {
// 		return nil
// 	}
// 	return sw.wait()
// }
