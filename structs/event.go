package structs

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

/*
eventBus := GenericEventBus{}
eventBus.Register(A, func(ctx context.Context, v T) {})
eventBus.Register(B, func(ctx context.Context, v T) {})
eventBus.Register(C, func(ctx context.Context, v T) {})
eventBus.Register(D, func(ctx context.Context, v T) {})

raw := json.Unmarshal()

event := Event[B]
eventBus.Handle(event)

*/
/*
type Callback[T any] func(ctx context.Context, v T) error

type EventBus[T any] struct {
	types map[T][]Callback
}

func (e EventBus) Handle(e T) {

}

*/

constraints.

type EventBus[T any] struct {
}

func (e EventBus[T]) Dispatch(ee Event[T]) {
	fmt.Println(ee)
}

type Event[T any] struct {
	name string
	t    T
}

func (e *Event[T]) T() T {
	return e.t
}

func (e *Event[T]) Hash() []byte {
	t := reflect.TypeOf(e.t)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("not struct")
	}

	fmt.Println(t.Kind())

	return nil
}

var ErrNotStruct = errors.New("only struct allowed")

func Hash(any interface{}) ([]byte, error) {
	t := reflect.TypeOf(any)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	fields := structFields(t, "")

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	for _, i := range fields {
		fmt.Println(i)
	}

	return nil, nil
}

func structFields(st reflect.Type, prefix string) []reflect.StructField {
	fields := make([]reflect.StructField, 0, st.NumField())

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		switch field.Type.Kind() {
		case reflect.Struct:
			fields = append(fields, structFields(field.Type, fmt.Sprintf("%s%s_", prefix, field.Name))...)
		default:
			field.Name = fmt.Sprintf("%s%s", prefix, field.Name)
			fields = append(fields, field)
		}
	}

	return fields
}
