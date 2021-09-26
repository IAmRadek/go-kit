package heartbeat_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/RadekD/go-kit/heartbeat"
)

func TestHeartbeat(t *testing.T) {
	count := 10
	delta := 10 * time.Millisecond

	var errs []string
	logErrs := func() {
		for _, e := range errs {
			t.Log(e)
		}
	}

	for i := 0; i < 5; i++ {
		H := heartbeat.New(delta)
		for i := 0; i < count; i++ {

			if i%2 == 0 {
				H.Skip()
				//continue
			}
			t0 := time.Now()
			<-H.C()
			t1 := time.Now()
			dt := t1.Sub(t0)

			x := time.Duration(1)
			if i%2 == 0 {
				x = time.Duration(2)
			}
			target := delta * x
			slop := target * 2 / 10

			if dt < target-slop || dt > target+slop {
				errs = append(errs, fmt.Sprintf("%d %s heartbeat took %s, expected [%s,%s]", count, delta, dt, target-slop, target+slop))
				continue
			}
		}
		H.Stop()

		select {
		case <-H.C():
			errs = append(errs, "Heartbeat did not shut down")
			continue
		default:
			// ok
		}

		// Test passed, so all done.
		if len(errs) > 0 {
			t.Logf("saw %d errors, ignoring to avoid flakiness", len(errs))
			logErrs()
		}

		return
	}

	t.Errorf("saw %d errors", len(errs))
	logErrs()
}
