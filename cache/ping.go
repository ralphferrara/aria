package cache

import "fmt"

//||------------------------------------------------------------------------------------------------||
//|| Redis/KeyDB: Ping
//||------------------------------------------------------------------------------------------------||

func (c *RedisCacheWrapper) Ping() error {
	if c == nil || c.Client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return c.Client.Ping(c.Ctx).Err()
}

//||------------------------------------------------------------------------------------------------||
//|| Memcached: Ping
//||------------------------------------------------------------------------------------------------||

func (c *MemcachedCacheWrapper) Ping() error {
	if c == nil || c.Client == nil {
		return fmt.Errorf("memcached client not initialized")
	}
	return c.Client.Ping()
}

//||------------------------------------------------------------------------------------------------||
//|| Memory: Ping (always ok)
//||------------------------------------------------------------------------------------------------||

func (c *MemoryCacheWrapper) Ping() error {
	if c == nil {
		return fmt.Errorf("memory cache not initialized")
	}
	return nil
}
