package ws_test

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/RadekD/go-kit/ws"
	"github.com/posener/wstest"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type testStruct struct {
	Test string
}
type testStructBad struct {
	Test int
}

type testStruct2 struct {
	Test *testStruct
}

func wsTestNoParams(r *ws.Request) {
	r.Respond("OK")
}

func wsTestString(r *ws.Request, test string) {
	if test == "test" {
		r.Respond("OK")
	} else {
		r.Respond("NOK")
	}
}

func wsTestStruct(r *ws.Request, test testStruct) {
	if test.Test == "test" {
		r.Respond("OK")
	} else {
		r.Respond("NOK")
	}
}
func wsTestPointer(r *ws.Request, test *testStruct) {
	if test.Test == "test" {
		r.Respond("OK")
	} else {
		r.Respond("NOK")
	}
}

func wsTestInsidePointer(r *ws.Request, test testStruct2) {
	if test.Test.Test == "test" {
		r.Respond("OK")
	} else {
		r.Respond("NOK")
	}
}

type packet struct {
	Id    int64
	Topic string
	Data  interface{}
}

type expected struct {
	packet
	Data json.RawMessage
}

func Test(t *testing.T) {
	log.Out = io.MultiWriter()

	handlers := []struct {
		Name    string
		Handler interface{}
	}{
		{"testNoParams", wsTestNoParams},
		{"testString", wsTestString},
		{"testStruct", wsTestStruct},
		{"testPointer", wsTestPointer},
		{"testInsidePointer", wsTestInsidePointer},
	}

	ws := &ws.WS{
		Log: log,
	}
	for _, h := range handlers {
		ws.RegisterHandler(h.Name, h.Handler)
	}

	tests := []struct {
		Test packet
	}{
		{packet{1, "testNoParams", nil}},
		{packet{2, "testString", "test"}},
		{packet{3, "testStruct", testStruct{"test"}}},
		{packet{4, "testPointer", &testStruct{"test"}}},
		{packet{5, "testInsidePointer", testStruct2{&testStruct{"test"}}}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Test.Topic, func(t *testing.T) {
			dial := wstest.NewDialer(ws, nil)
			conn, _, err := dial.Dial("ws://example.org/websocket", nil)
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()

			conn.WriteJSON(test.Test)

			var resp packet
			conn.ReadJSON(&resp)

			if resp.Data != "OK" {
				t.Fail()
			}
		})
	}

	t.Run("testBadType", func(t *testing.T) {
		dial := wstest.NewDialer(ws, nil)
		conn, _, err := dial.Dial("ws://example.org/websocket", nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		conn.WriteJSON(packet{1, "testString", 1})

		var resp packet
		err = conn.ReadJSON(&resp)
		if err == nil {
			t.Fail()
		}
	})

	t.Run("testBadType2", func(t *testing.T) {
		dial := wstest.NewDialer(ws, nil)
		conn, _, err := dial.Dial("ws://example.org/websocket", nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		var resp packet
		conn.WriteJSON(packet{1, "testStruct", testStructBad{1}})
		err = conn.ReadJSON(&resp)
		if err == nil {
			t.Fail()
		}
	})

}
