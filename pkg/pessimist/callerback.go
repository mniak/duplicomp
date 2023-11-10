package pessimist

type Callerback struct {
	succeded         bool
	successCallbacks []func()
	failureCallbacks []func()
}

func (self *Callerback) Callback() {
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

func (self *Callerback) OnFailure(fn func()) {
	self.failureCallbacks = append(self.failureCallbacks, fn)
}

func (self *Callerback) OnSuccess(fn func()) {
	self.successCallbacks = append(self.successCallbacks, fn)
}

func (self *Callerback) Succeeded() {
	self.succeded = true
}
