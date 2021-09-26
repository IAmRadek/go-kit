//Package ping provides hook for pinging websocket every 5 seconds
//
//TODO:
//	Add configuration
package ping

import (
	"time"

	"github.com/RadekD/go-kit/ws"
)

//Hook is
func Hook(c *ws.Connection) {
	ctx := c.Context()
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.Send(0, "pong", time.Now().UnixNano())
			}
		}
	}()
}
