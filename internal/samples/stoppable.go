package samples

type Stoppable interface {
	Stop()
}

type stoppable func()

func (s stoppable) Stop() {
	s()
}
