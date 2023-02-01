package cache

import "time"

type RefreshableCache[V interface{}, K string | int | uint] struct {
	Cache[V, K]
	refresh func(K) *V
}

func InitRefreshCache[V interface{}, K string | int | uint](ttl time.Duration, refreshFunc func(K) *V) RefreshableCache[V, K] {
	rc := RefreshableCache[V, K]{
		Cache:   InitCache[V, K](ttl),
		refresh: refreshFunc,
	}

	return rc
}

func (rc *RefreshableCache[T, K]) GetEntry(key K) (*T, bool) {
	entry, exists := rc.Cache.GetEntry(key)
	if !exists {
		entry = rc.refresh(key)
		exists = entry == nil
	}
	return entry, exists
}
