package cancel

import "time"

var nocancelChan = make(chan struct{})

// NeverCancel returns a Canceller which never cancels.
func NeverCancel() Canceller {
	return nevercancel{}
}

var _ Canceller = nevercancel{}

type nevercancel struct{}

func (nevercancel) Deadline() (time.Time, bool) {
	return time.Time{}, false
}
func (nevercancel) Done() <-chan struct{} {
	return nocancelChan
}
func (nevercancel) Err() error {
	return nil
}
