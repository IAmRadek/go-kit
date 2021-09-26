package online

import (
	"expvar"

	"github.com/RadekD/go-kit/ws"
)

var online = expvar.NewInt("websocket_connections")

//PreHook increments online
func PreHook(conn *ws.Connection) {
	online.Add(1)
}

//PostHook decrements online
func PostHook(conn *ws.Connection) {
	online.Add(-1)
}

//Get returns online number in safe way
func Get() int64 {
	return online.Value()
}
