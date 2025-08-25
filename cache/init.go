//||------------------------------------------------------------------------------------------------||
//|| Cache Package: Initialization
//|| init.go
//||------------------------------------------------------------------------------------------------||

package cache

import (
	"fmt"

	"github.com/ralphferrara/aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| Cache: Globals (Redis, KeyDB, Memcached, Memory)
//||------------------------------------------------------------------------------------------------||

var (
	Redis     = map[string]*RedisCacheWrapper{}
	KeyDB     = map[string]*RedisCacheWrapper{}
	Memcached = map[string]*MemcachedCacheWrapper{}
	Memory    = map[string]*MemoryCacheWrapper{}
)

//||------------------------------------------------------------------------------------------------||
//|| Cache: Init - Connects all caches from config
//||------------------------------------------------------------------------------------------------||

func Init() error {
	cfg := config.GetConfig()
	for name, cacheCfg := range cfg.Cache {
		switch cacheCfg.Backend {
		case "redis":
			client, ctx, err := connectRedis(cacheCfg)
			if err != nil {
				return fmt.Errorf("failed to connect to redis '%s': %w", name, err)
			}
			Redis[name] = &RedisCacheWrapper{
				Name:   name,
				Client: client,
				Ctx:    ctx,
			}
		case "keydb":
			client, ctx, err := connectRedis(cacheCfg)
			if err != nil {
				return fmt.Errorf("failed to connect to keydb '%s': %w", name, err)
			}
			KeyDB[name] = &RedisCacheWrapper{
				Name:   name,
				Client: client,
				Ctx:    ctx,
			}
		case "memcached":
			client, err := connectMemcached(cacheCfg)
			if err != nil {
				return fmt.Errorf("failed to connect to memcached '%s': %w", name, err)
			}
			Memcached[name] = &MemcachedCacheWrapper{
				Name:   name,
				Client: client,
			}
		case "memory":
			Memory[name] = &MemoryCacheWrapper{
				Name:  name,
				Store: initMemoryStore(),
			}
		default:
			return fmt.Errorf("unsupported cache backend: %s", cacheCfg.Backend)
		}
	}
	return nil
}
