package cancel

// WithCanceller returns a Canceller which will cancel when the
// returned function is executed. If `canc` gets cancelled, the returned
// Canceller will also be cancelled.
func WithCanceller(canc Canceller) (Canceller, func()) {

	t, ok := canc.Deadline()

	return newCanceller(canc, t, ok, false)
}
