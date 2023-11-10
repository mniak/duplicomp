package ps121

type PessimisticCallerback struct {
	succeded         bool
	successCallbacks []func()
	failureCallbacks []func()
}

func (self *PessimisticCallerback) Callback() {
	if self.succeded {
		for _, dispose := range self.successCallbacks {
			if dispose != nil {
				dispose()
			}
		}
	} else {
		for _, dispose := range self.failureCallbacks {
			if dispose != nil {
				dispose()
			}
		}
	}
}

func (self *PessimisticCallerback) OnFailure(fn func()) {
	self.failureCallbacks = append(self.failureCallbacks, fn)
}

func (self *PessimisticCallerback) OnSuccess(fn func()) {
	self.successCallbacks = append(self.successCallbacks, fn)
}

func (self *PessimisticCallerback) Succeeded() {
	self.succeded = true
}
