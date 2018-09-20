package cancel // import "github.com/deanveloper/cancel"

import (
	"errors"
	"sync"
	"time"
)

// ErrCancel is an error returned if the task was
// cancelled using the Cancel() function
var ErrCancel = errors.New("task cancelled")

// ErrDeadline is an error returned if the task was
// cancelled because it has reached its deadline.
var ErrDeadline = errors.New("deadline reached")

// Canceller is similar to a context.Context, but is strictly
// used for cancelling tasks safely, and nothing else.
//
// Cancelling an already-cancelled task results in a no-op.
type Canceller interface {

	// Deadline tells you the deadline for this canceller. `ok` is true
	// if a deadline for this canceller exists, even if the deadline has
	// already passed.
	Deadline() (deadline time.Time, ok bool)

	// Done returns a channel which is closed when the
	// task has been cancelled.
	Done() <-chan struct{}

	// Err returns nil if the task has not been cancelled.
	// otherwise, it returns ErrCancelled or ErrDeadline.
	Err() error
}

// canc - the parent canceller
// t    - the deadline for this handler
// hasT - if this _has_ a time to cancel at
// useT - if it should actually use the deadline.
//        we do not want to listen to the deadline if
//        the parent is going to expire at the same time anyway.
func newCanceller(canc Canceller, t time.Time, hasT, useT bool) (Canceller, func()) {

	intern := &canceller{
		t:      t,
		hasT:   hasT,
		ch:     make(chan struct{}),
		parent: canc,
	}

	// close when the parent closes.
	// creates a goroutine which ends
	// when retval gets closed.
	go func() {
		if useT {
			select {
			case <-time.After(time.Until(t)):
				intern.close(ErrDeadline)
			case <-intern.parent.Done():
				intern.close(canc.Err())
			case <-intern.Done():
			}
		} else {
			select {
			case <-intern.parent.Done():
				intern.close(canc.Err())
			case <-intern.Done():
			}
		}
	}()

	return intern, func() { intern.close(ErrCancel) }
}

type canceller struct {
	parent Canceller
	t      time.Time
	hasT   bool
	ch     chan struct{}
	err    error
	once   sync.Once
}

// assert that canceller implements our interface Canceller
var _ Canceller = &canceller{}

func (d *canceller) Deadline() (time.Time, bool) {
	return d.t, d.hasT
}
func (d *canceller) Done() <-chan struct{} {
	return d.ch
}
func (d *canceller) Err() error {
	return d.err
}
func (d *canceller) close(err error) {
	d.once.Do(func() {
		d.err = err
		close(d.ch)
	})
}
