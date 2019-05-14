package yamux

import (
	"sync"
	"time"
)

type deadline struct {
	C     <-chan time.Time
	timer *time.Timer
	mu    sync.Mutex
}

func newDeadline() *deadline {
	d := new(deadline)
	d.timer = time.NewTimer(0)
	d.C = d.timer.C
	d.set(time.Time{})
	return d
}

// setExpired forces the deadline to expire now.
func (d *deadline) setExpired() {
	d.mu.Lock()
	if d.timer != nil {
		d.timer.Reset(0)
	}
	d.mu.Unlock()
}

// set sets the deadline.
func (d *deadline) set(t time.Time) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer == nil {
		return ErrStreamClosed
	}
	d.timer.Stop()
	select {
	case <-d.timer.C:
	default:
	}
	if !t.IsZero() {
		// negative timeout triggers immediately.
		d.timer.Reset(time.Until(t))
	}
	return nil
}

func (d *deadline) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer == nil {
		return
	}
	d.timer.Stop()
	select {
	case <-d.timer.C:
	default:
	}
	d.timer = nil
}
