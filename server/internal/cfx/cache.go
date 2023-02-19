package cfx

import (
	cfx_models "degrens/panel/internal/db/models/cfx"
	"sync"
	"time"
)

type Cache struct {
	Mutex   sync.Mutex
	Players PlayersCache
}

var cache Cache

func init() {
	cache = Cache{
		Mutex: sync.Mutex{},
		Players: PlayersCache{
			Data:      []cfx_models.User{},
			UpdatedAt: time.UnixMilli(1),
		},
	}
}

func getCache() *Cache {
	return &cache
}

func UnlockCache() {
	cache.Mutex.TryLock()
	cache.Mutex.Unlock()
}
