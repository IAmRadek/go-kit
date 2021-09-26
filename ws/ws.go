package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var ErrUnknowHandler = errors.New("ws: unknown handler")
var ErrInvalidPacket = errors.New("ws: invalid packet")

//SM is a helper for simple messages with success and message only

type sm struct {
	Success bool
	Message string
}

func SM(success bool, message string) sm {
	return sm{success, message}
}

type handler func(r *Request) error

//Hook is function called before main loop or after connection was closed
type Hook func(conn *Connection)

//Connection holds all data and websocket connection
type Connection struct {
	l    sync.Mutex
	conn *websocket.Conn
	ctx  context.Context

	Request *http.Request
	Values  sync.Map
}

func (conn *Connection) Context() context.Context {
	return conn.ctx
}

func (conn *Connection) Send(id int64, topic string, data interface{}) error {
	conn.l.Lock()
	defer conn.l.Unlock()
	return conn.conn.WriteJSON(struct {
		ID    int64
		Topic string
		Data  interface{}
	}{id, topic, data})
}

//Request represents data sends from client to server
type Request struct {
	C *Connection `json:"-"`

	ID    int64
	Topic string
	Data  json.RawMessage
}

//Respond sends response to client and returns error if failed
func (r *Request) Respond(data interface{}) error {
	return r.C.Send(r.ID, r.Topic, data)
}

type WS struct {
	ErrorHandler func(c *Connection, err error)
	Upgrader     websocket.Upgrader

	handlers  map[string]handler
	preHooks  []Hook
	postHooks []Hook
}

var rTest = reflect.TypeOf(&Request{})

//RegisterHandler register function for entered topic
// note that handlerFn is interface{} but should be function (func(r *Request, optionalParameter string))
func (m *WS) RegisterHandler(topic string, handlerFn interface{}) {
	if m.handlers == nil {
		m.handlers = make(map[string]handler)
	}
	if _, ok := m.handlers[topic]; ok {
		panic("topic already registered")
	}

	fn := reflect.TypeOf(handlerFn)
	if fn.Kind() != reflect.Func || fn.NumIn() < 1 || fn.NumIn() > 2 {
		panic("handler not function")
	}

	arg1 := fn.In(0)
	if arg1 != rTest {
		panic("handler is invalid")
	}

	var in []reflect.Type
	in = append(in, arg1)

	if fn.NumIn() > 1 {
		arg2 := fn.In(1)
		in = append(in, arg2)
	}

	m.handlers[topic] = func(r *Request) error {

		var toFn []reflect.Value
		toFn = append(toFn, reflect.ValueOf(r))

		if fn.NumIn() > 1 {
			var dataValue = reflect.New(in[1])

			err := json.Unmarshal(r.Data, dataValue.Interface())
			if err != nil {
				return err
			}

			toFn = append(toFn, dataValue.Elem())
		}
		reflect.ValueOf(handlerFn).Call(toFn)

		return nil
	}
}

//AddPreHook appends pre main loop hook
func (m *WS) AddPreHook(hook Hook) {
	m.preHooks = append(m.preHooks, hook)
}

//AddPostHook appends post request hook
func (m *WS) AddPostHook(hook Hook) {
	m.postHooks = append(m.postHooks, hook)
}

func (m *WS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws, err := m.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.ErrorHandler(nil, fmt.Errorf("Error: %w = could not upgrade connection", err))
		return
	}
	defer ws.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	conn := &Connection{
		conn: ws,
		ctx:  ctx,

		Request: r,
	}

	for _, hook := range m.preHooks {
		hook(conn)
	}

	for {
		var message Request
		err := ws.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				m.ErrorHandler(conn, fmt.Errorf("Error: %w = socket closed", err))
			}
			break
		}

		if message.ID == 0 || message.Topic == "" {
			m.ErrorHandler(conn, ErrInvalidPacket)
			break
		}

		handler, ok := m.handlers[message.Topic]
		if !ok {
			m.ErrorHandler(conn, errors.Wrapf(ErrUnknowHandler, "handler: %s", message.Topic))
			continue
		}

		message.C = conn

		err = handler(&message) // only error I can get here is from json.Unmarshal
		if err != nil {
			m.ErrorHandler(conn, errors.Wrapf(err, "handler: %s", message.Topic))
			break
		}
	}

	for _, hook := range m.postHooks {
		hook(conn)
	}
}

/*
&WS{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  10240,
			WriteBufferSize: 10240,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
*/
