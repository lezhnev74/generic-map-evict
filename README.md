# Generic Map With Eviction

This is a pool that keeps objects as K->V. We can ask for cached objects.
If the object is not in the cache it will be created with a factory function and put to the cache.
The cache will evict items with respect to max allowed capacity.
Upon requesting N values it returns N values even if N exceeds pool capacity.

## API

```go
capacity := 5
factoryFn := func (key Comparable) any {
return NewObject(key);
}

// 0. Make a new MapCache
m := NewMapCacheEvictUnused(capacity, factoryFn)

// 1. Request new values - will call factory 3 times
instances := m.Get([]string{"a1", "a2", "a3"})
// len(instances) == 3
// m.Count() == 3

// 2. Request some existing and new ones (factory is called once):
instances := m.Get([]string{"a1", "a4"})
// len(instances) == 2
// m.Count() == 4

// 3. Request the number of values(6) that exceeds capacity(5), factory called 2 times
instances := m.Get([]string{"a1", "a2", "a3", "a4", "a5", "a6"})
// len(instances) == 6
// m.Count() == 5

```

Optionally, you can provide a destructor for evicted items:
```go
m := NewMapCacheEvictUnused(capacity, factoryFn)
m.SetFreeFn(func (key string, val any) { close(val) }) // <-- each evicted item goes through this func
```