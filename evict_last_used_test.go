package mapCache

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestEvictionLastUsed(t *testing.T) {
	t.Run("test it never exceeds capacity", func(t *testing.T) {
		type test struct {
			cap          int
			initialState []string
			getItems     []string
		}

		tests := []test{
			{ // evict 0: add to empty
				1,
				[]string{},
				[]string{"a"},
			},
			{ // evict 0: add to half empty
				2,
				[]string{"a"},
				[]string{"b"},
			},
			{ // evict 1
				1,
				[]string{"a"},
				[]string{"b"},
			},
			{ // evict many
				1,
				[]string{"a"},
				[]string{"b", "c"},
			},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
				m := NewMapCacheEvictUnused(tt.cap, func(key string) any { return rand.Int() })
				m.Get(tt.initialState)
				time.Sleep(time.Millisecond) // make sure new items are younger
				m.Get(tt.getItems)

				require.True(t, m.Count() <= m.cap)
			})
		}
	})
}
