package cancel_test

import (
	"testing"
	"time"

	"github.com/deanveloper/cancel"
)

const interval = 10 * time.Millisecond

// a task that should take ~100ms to complete but is cancellable every 10ms
func cancellableTask(canc cancel.Canceller, ch chan<- int) {
	for i := 0; i < 10; i++ {
		time.Sleep(interval)
		select {
		case <-canc.Done():

		default:
		}
		ch <- i
	}
}
func TestNeverCancel(t *testing.T) {
	ch := make(chan int)
	canc := cancel.NeverCancel()
	go cancellableTask(canc, ch)

	for i := 0; i < 10; i++ {
		select {
		case <-time.After(interval + 3*time.Millisecond):
			t.Fatalf("timed out: i(%d)", i)
		case <-canc.Done():
			t.Fatalf("task cancelled: i(%d) err(%s)", i, canc.Err())
		case <-ch:
		}
	}

	if canc.Err() != nil {
		t.Fatalf("task err not nil: %v", canc.Err())
	}
	if dl, ok := canc.Deadline(); ok {
		t.Fatalf("task has deadline: %s", dl)
	}
	select {
	case <-canc.Done():
		t.Fatalf("task cancelled (later)")
	case <-time.After(50 * time.Millisecond):
		break
	}
}

func TestDeadline(t *testing.T) {
	ch := make(chan int)
	// cancel after the 3rd iteration
	canc, f := cancel.WithTimeLimit(cancel.NeverCancel(), 3*interval+2*time.Millisecond)
	defer f()

	go cancellableTask(canc, ch)

	for i := 0; i < 3; i++ {
		select {
		case <-time.After(interval + 3*time.Millisecond):
			t.Errorf("timed out: i(%d)", i)
		case <-canc.Done():
			t.Errorf("task cancelled early: i(%d) err(%s)", i, canc.Err())
		case <-ch:
		}
	}

	// special case on the 4th (i=3) iteration
	select {
	case <-time.After(interval + 7*time.Millisecond):
		t.Errorf("timed out: i(%d)", 3)
	case <-canc.Done():
		break // yay!
	case i := <-ch:
		t.Errorf("task continued i(%d)", i)
	}

	if canc.Err() == nil {
		t.Errorf("task err not set!")
	}
	if _, ok := canc.Deadline(); !ok {
		t.Errorf("task does not have deadline!")
	}
}

func TestWithCancelFunc(t *testing.T) {
	ch := make(chan int)
	// cancel after the 3rd iteration
	canc, f := cancel.WithCanceller(cancel.NeverCancel())
	defer f() // safe to do as multiple calls to f() will no-op
	go cancellableTask(canc, ch)

	time.AfterFunc(3*interval+3*time.Millisecond, f)

	for i := 0; i < 3; i++ {
		select {
		case <-time.After(interval + 3*time.Millisecond):
			t.Errorf("timed out: i(%d)", i)
		case <-canc.Done():
			t.Errorf("task cancelled early: i(%d) err(%s)", i, canc.Err())
		case <-ch:
		}
	}

	// special case on the 4th (i=3) iteration
	select {
	case <-time.After(interval + 7*time.Millisecond):
		t.Errorf("timed out: i(%d)", 3)
	case <-canc.Done():
		break // yay!
	case i := <-ch:
		t.Errorf("task continued i(%d)", i)
	}

	if canc.Err() == nil {
		t.Errorf("task err not set!")
	}
	if dl, ok := canc.Deadline(); ok {
		t.Errorf("task has a deadline: %s", dl)
	}
}

func TestParentCancel(t *testing.T) {
	ch := make(chan int)

	// cancel after the 3rd iteration
	root, f := cancel.WithCanceller(cancel.NeverCancel())
	time.AfterFunc(3*interval+2*time.Millisecond, f)

	// cancels itself _around_ the same time as the root, but a little
	// bit later.
	parent, _ := cancel.WithTimeLimit(root, 3*interval+5*time.Millisecond)

	// child canceller which cancels after the 5th iteration (but will really
	// cancel after the 3rd iteration)
	canc, _ := cancel.WithTimeLimit(parent, 5*interval+time.Millisecond)

	dl, _ := canc.Deadline()
	if dl.Truncate(time.Millisecond) != time.Now().Add(3*interval+5*time.Millisecond).Truncate(time.Millisecond) {
		t.Errorf("incorrect deadline: %s", time.Until(dl))
	}

	go cancellableTask(canc, ch)

	for i := 0; i < 3; i++ {
		select {
		case <-time.After(interval + 3*time.Millisecond):
			t.Errorf("timed out: i(%d)", i)
		case <-canc.Done():
			t.Errorf("task cancelled early: i(%d) err(%s)", i, canc.Err())
		case <-ch:
		}
	}

	// special case on the 4th (i=3) iteration
	select {
	case <-time.After(interval + 7*time.Millisecond):
		t.Errorf("timed out: i(%d)", 3)
	case <-canc.Done():
		break // yay!
	case i := <-ch:
		t.Errorf("task continued i(%d)", i)
	}

	if canc.Err() == nil {
		t.Errorf("task err not set!")
	}
	if _, ok := canc.Deadline(); !ok {
		t.Errorf("task does not have a deadline: %s", dl)
	}
}
