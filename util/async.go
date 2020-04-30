package util

import "git.devops.com/go/odm/types"

// FinalAction is something already done.
type finalAction struct {
	err error
}

// WaitDone implement
func (r *finalAction) WaitDone() error {
	return r.err
}

type waitingAction struct {
	err      error
	doneChan chan bool
}

func (a *waitingAction) WaitDone() error {
	<-a.doneChan
	return a.err
}

// AsyncError return an asyncAction
func AsyncError(err error) types.AsyncAction {
	a := new(finalAction)
	a.err = err
	return a
}

// RunAsync make async call easier
func RunAsync(f func() error) types.AsyncAction {
	action := new(waitingAction)
	action.doneChan = make(chan bool, 1)
	go func() {
		action.err = f()
		action.doneChan <- true
	}()
	return action
}
