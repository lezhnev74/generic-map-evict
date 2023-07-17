package mapCache

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	capacity := 5
	factoryFn := func(key string) any {
		return key
	}

	// 0. Make a new MapCache
	m := NewMapCacheEvictUnused(capacity, factoryFn)

	//1. Request new values - will call Factory 3 times
	instances := m.Get([]string{"a1", "a2", "a3"})
	require.Len(t, instances, 3)
	require.Equal(t, 3, m.Count())

	// 2. Request some existing and new ones (Factory is called once):
	instances = m.Get([]string{"a1", "a4"})
	require.Len(t, instances, 2)
	require.Equal(t, 4, m.Count())

	// 3. Request the number of values(6) that exceeds capacity(5), Factory called 2 times
	instances = m.Get([]string{"a1", "a2", "a3", "a4", "a5", "a6"})
	require.Len(t, instances, 6)
	require.Equal(t, 5, m.Count())
}

func TestFreeProcedure(t *testing.T) {
	capacity := 1
	factoryFn := func(key string) any { return key }
	freed := []string{}
	m := NewMapCacheEvictUnused(capacity, factoryFn)
	m.SetFreeFn(func(key string, val any) {
		freed = append(freed, key)
	})

	m.Get([]string{"a"})
	time.Sleep(time.Millisecond)
	m.Get([]string{"b"})

	require.EqualValues(t, []string{"a"}, freed)
}
