package mapCache

type Getter[T comparable] func(key []T) map[T]any

type MapCache[T comparable] struct {
	cap       int
	m         map[T]any
	factoryFn func(key T) any
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
	m.evict(keys)
	return ret
}
func (m *MapCache[T]) Count() int { return len(m.m) }

// EvictionPolicy gets currently requested keys and decides if other keys should retire
type EvictionPolicy[T comparable] func(keys []T)

func NewMapCacheEvictUnused[T comparable](capacity int, fn func(key T) any) MapCache[T] {
	var m MapCache[T]
	m = MapCache[T]{
		cap:       capacity,
		m:         make(map[T]any, capacity),
		factoryFn: fn,
		evict:     newEvictFirstUsed(&m),
	}
	return m
}
