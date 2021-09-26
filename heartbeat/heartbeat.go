/*
Package heartbeat is a extension to time.Ticker that privides skipping capability.
I think is useful, when you want to perform action (eg. closeing some connections), but only after some other action did not occur in some time,
but if it did, you want to extend the waiting time.
*/
package heartbeat

import "time"

//Heartbeat interface highlights functionality of this package
type Heartbeat interface {
	//C() is channel on which you should listen to
	C() <-chan time.Time
	//Skip allows you to skip a heat beat
	Skip()
	//Stop allows you to stop a heart
	Stop()
}

type heartbeat struct {
	c       chan time.Time
	stopped bool
	skipped bool
}

func (h *heartbeat) C() <-chan time.Time {
	return h.c
}

func (h *heartbeat) Skip() {
	h.skipped = true
}
func (h *heartbeat) Stop() {
	h.stopped = true
}

//New creates new heartbeat interface
func New(d time.Duration) Heartbeat {
	t := time.NewTicker(d)

	c := make(chan time.Time, 1)
	h := heartbeat{
		c:       c,
		stopped: false,
		skipped: false,
	}
	go func() {
		for !h.stopped {
			select {
			case <-t.C:
				{
					if h.skipped {
						h.skipped = false
						continue
					}
					c <- time.Now()
				}
			}
		}
		t.Stop()
	}()
	return &h
}
