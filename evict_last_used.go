package mapCache

import (
	"golang.org/x/exp/slices"
	"time"
)

type usedKey[T comparable] struct {
	key  T
	used time.Time
}

func newEvictFirstUsed[T comparable](m *MapCache[T]) EvictionPolicy[T] {
	usedKeys := make(map[T]time.Time, m.cap) // last usage map, retire oldest

	return func(keys []T) {
		now := time.Now()
		for _, key := range keys {
			usedKeys[key] = now
		}
		// make a slice of all key sorted by last used
		allKeysSorted := make([]usedKey[T], 0, len(usedKeys))
		for key, keyUsed := range usedKeys {
			allKeysSorted = append(allKeysSorted, usedKey[T]{key, keyUsed})
		}
		slices.SortFunc(allKeysSorted, func(a, b usedKey[T]) bool { return a.used.Before(b.used) })
		retireCount := len(usedKeys) - m.cap
		for i := 0; i < retireCount; i++ {
			retireKey := allKeysSorted[i].key
			delete(m.m, retireKey)
			delete(usedKeys, retireKey)
		}
	}
}
