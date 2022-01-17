package structs

import (
	"testing"
)

type B struct {
	A string
	B int
}

func Test_base(t *testing.T) {

	g := new(Event[B])

	eb := EventBus[B]{}

	eb.Dispatch(*g)
}
