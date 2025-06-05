package syncx

import "sync"

type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return value, ok
	}
	return v.(V), ok
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *Map[K, V]) Clear() {
	m.m.Clear()
}

func (m *Map[K, V]) CompareAndDelete(key, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

func (m *Map[K, V]) CompareAndSwap(key, old, new V) (swapped bool) {
	return m.m.CompareAndSwap(key, old, new)
}

func (m *Map[K, V]) Delete(key V) {
	m.m.Delete(key)
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, ok := m.m.LoadAndDelete(key)
	if !ok {
		return value, ok
	}
	return v.(V), ok
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, ok := m.m.LoadOrStore(key, value)
	if !ok {
		return value, ok
	}
	return v.(V), ok
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func (m *Map[K, V]) Swap(key K, value V) (previous any, loaded bool) {
	return m.m.Swap(key, value)
}
