package cancel

import (
	"time"
)

// WithDeadline returns a Canceller which will cancel after
// `t` has been reached. If `canc` gets cancelled, the returned
// Canceller will also be cancelled.
func WithDeadline(canc Canceller, t time.Time) (Canceller, func()) {

	deadline := t

	if oDead, ok := canc.Deadline(); ok {
		if oDead.Before(deadline) {
			deadline = oDead
		}
	}

	return newCanceller(canc, deadline, true, t == deadline)
}

// WithTimeLimit is the same as WithDeadline, but takes a `time.Duration`
// instead of a `time.Time`.
func WithTimeLimit(canc Canceller, d time.Duration) (Canceller, func()) {
	return WithDeadline(canc, time.Now().Add(d))
}
