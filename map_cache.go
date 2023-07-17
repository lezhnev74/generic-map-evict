package mapCache

type Getter[T comparable] func(key []T) map[T]any

type MapCache[T comparable] struct {
	cap       int
	m         map[T]any
	factoryFn func(key T) any
	freeFn    func(key T, object any) // free up evicted data (optional)
	evict     EvictionPolicy[T]
}

func (m *MapCache[T]) Get(keys []T) map[T]any {
	ret := make(map[T]any, len(keys))
	for _, key := range keys {
		if _, ok := m.m[key]; !ok {
			m.m[key] = m.factoryFn(key)
		}
		ret[key] = m.m[key]
	}
	evictedKeys := m.evict(keys)
	for _, key := range evictedKeys {
		m.freeFn(key, m.m[key])
		delete(m.m, key)
	}
	return ret
}
func (m *MapCache[T]) Count() int                            { return len(m.m) }
func (m *MapCache[T]) SetFreeFn(freeFn func(key T, val any)) { m.freeFn = freeFn }

// EvictionPolicy gets currently requested keys and decides if other keys should retire
// returns keys that must retire
type EvictionPolicy[T comparable] func(keys []T) []T

func NewMapCacheEvictUnused[T comparable](capacity int, fn func(key T) any) MapCache[T] {
	var m MapCache[T]
	m = MapCache[T]{
		cap:       capacity,
		m:         make(map[T]any, capacity),
		factoryFn: fn,
		freeFn:    func(key T, object any) {}, // no-op by default
		evict:     newEvictFirstUsed(&m),
	}
	return m
}
