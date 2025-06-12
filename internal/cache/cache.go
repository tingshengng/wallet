package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type InMemoryCache struct {
	cache *cache.Cache
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: cache.New(5*time.Minute, 10*time.Minute), // Default cleanup interval and item expiration
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	value, found := c.cache.Get(key)
	if !found {
		return nil, false
	}
	return value, true
}

func (c *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

func (c *InMemoryCache) Delete(key string) {
	c.cache.Delete(key)
}
