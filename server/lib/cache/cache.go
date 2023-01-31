package cache

import "time"

type Cache[T interface{}, K string | int | uint] struct {
	data map[K]CacheEntry[T]
	ttl  time.Duration
}

type CacheEntry[T interface{}] struct {
	data        T
	lastUpdated time.Time
}

func InitCache[T interface{}, K string | int | uint](ttl time.Duration) Cache[T, K] {
	cache := *&Cache[T, K]{
		data: map[K]CacheEntry[T]{},
		ttl:  ttl,
	}
	return cache
}

func (c *Cache[T, K]) Has(key K) bool {
	_, ok := c.data[key]
	return ok
}

func (c *Cache[T, K]) AddEntry(key K, data T) {
	entry := CacheEntry[T]{
		data:        data,
		lastUpdated: time.Now(),
	}
	c.data[key] = entry
}

func (c *Cache[T, K]) GetEntry(key K) (*T, bool) {
	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}
	if entry.lastUpdated.Add(c.ttl).Before(time.Now()) {
		return nil, false
	}
	return &entry.data, true
}
