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
	accessMap := make(map[T]time.Time, m.cap) // last usage map, retire oldest

	return func(keys []T) []T {
		now := time.Now()
		for _, key := range keys {
			accessMap[key] = now
		}
		// make a slice of all key sorted by last used
		allKeysSorted := make([]usedKey[T], 0, len(m.m))
		for key, _ := range m.m {
			lastUsed := accessMap[key]
			allKeysSorted = append(allKeysSorted, usedKey[T]{key, lastUsed})
		}
		slices.SortFunc(allKeysSorted, func(a, b usedKey[T]) bool { return a.used.Before(b.used) })
		retireCount := len(accessMap) - m.cap
		if retireCount <= 0 {
			return nil
		}
		retireKeys := make([]T, 0, retireCount)
		for i := 0; i < retireCount; i++ {
			retireKey := allKeysSorted[i].key
			retireKeys = append(retireKeys, retireKey)
			delete(accessMap, retireKey)
		}
		return retireKeys
	}
}
