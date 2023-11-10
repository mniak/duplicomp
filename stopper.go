package ps121

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
