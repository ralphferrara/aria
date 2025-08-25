//||------------------------------------------------------------------------------------------------||
//|| Cache Package: API
//|| api.go
//||------------------------------------------------------------------------------------------------||

package cache

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

//||------------------------------------------------------------------------------------------------||
//|| Redis/KeyDB: Set Key
//||------------------------------------------------------------------------------------------------||

func (c *RedisCacheWrapper) Set(key string, value interface{}, ttl time.Duration) error {
	return c.Client.Set(c.Ctx, key, value, ttl).Err()
}

//||------------------------------------------------------------------------------------------------||
//|| Redis/KeyDB: Get Key
//||------------------------------------------------------------------------------------------------||

func (c *RedisCacheWrapper) Get(key string) (string, error) {
	return c.Client.Get(c.Ctx, key).Result()
}

//||------------------------------------------------------------------------------------------------||
//|| Memcached: Set Key
//||------------------------------------------------------------------------------------------------||

func (c *MemcachedCacheWrapper) Set(key string, value []byte, ttl int32) error {
	return c.Client.Set(&memcache.Item{Key: key, Value: value, Expiration: ttl})
}

//||------------------------------------------------------------------------------------------------||
//|| Memcached: Get Key
//||------------------------------------------------------------------------------------------------||

func (c *MemcachedCacheWrapper) Get(key string) ([]byte, error) {
	item, err := c.Client.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Memory: Set Key
//||------------------------------------------------------------------------------------------------||

func (c *MemoryCacheWrapper) Set(key string, value interface{}) {
	c.Store[key] = value
}

//||------------------------------------------------------------------------------------------------||
//|| Memory: Get Key
//||------------------------------------------------------------------------------------------------||

func (c *MemoryCacheWrapper) Get(key string) (interface{}, bool) {
	v, ok := c.Store[key]
	return v, ok
}
