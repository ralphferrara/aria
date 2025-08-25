//||------------------------------------------------------------------------------------------------||
//|| Cache Package: Types & Configs
//|| types.go
//||------------------------------------------------------------------------------------------------||

package cache

import (
	"context"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
)

//||------------------------------------------------------------------------------------------------||
//|| Cache: Backend Enum
//||------------------------------------------------------------------------------------------------||

type CacheBackend string

const (
	BackendRedis     CacheBackend = "REDIS"
	BackendKeyDB     CacheBackend = "KEYDB"
	BackendMemcached CacheBackend = "MEMCACHED"
	BackendMemory    CacheBackend = "MEMORY"
)

//||------------------------------------------------------------------------------------------------||
//|| Redis/KeyDB Cache Wrapper
//||------------------------------------------------------------------------------------------------||

type RedisCacheWrapper struct {
	Name   string
	Client *redis.Client
	Ctx    context.Context
}

//||------------------------------------------------------------------------------------------------||
//|| Memcached Cache Wrapper
//||------------------------------------------------------------------------------------------------||

type MemcachedCacheWrapper struct {
	Name   string
	Client *memcache.Client
}

//||------------------------------------------------------------------------------------------------||
//|| Memory Cache Wrapper
//||------------------------------------------------------------------------------------------------||

type MemoryCacheWrapper struct {
	Name  string
	Store map[string]interface{}
}
