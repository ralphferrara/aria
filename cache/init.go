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
//|| Init (build all caches from main config)
//||------------------------------------------------------------------------------------------------||

func Init(cfg *config.Config) (
	map[string]*RedisCacheWrapper,
	map[string]*RedisCacheWrapper,
	map[string]*MemcachedCacheWrapper,
	map[string]*MemoryCacheWrapper,
	error,
) {

	//||------------------------------------------------------------------------------------------------||
	//|| Output Maps
	//||------------------------------------------------------------------------------------------------||

	redisMap := make(map[string]*RedisCacheWrapper)
	keydbMap := make(map[string]*RedisCacheWrapper)
	memcMap := make(map[string]*MemcachedCacheWrapper)
	memMap := make(map[string]*MemoryCacheWrapper)

	//||------------------------------------------------------------------------------------------------||
	//|| Loop Configured Caches
	//||------------------------------------------------------------------------------------------------||

	for name, c := range cfg.Cache {
		switch c.Backend {

		//||------------------------------------------------------------------------------------------------||
		//|| Redis
		//||------------------------------------------------------------------------------------------------||

		case "redis", "REDIS":
			client, ctx, err := connectRedis(c)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("cache '%s' redis connect failed: %w", name, err)
			}
			fmt.Printf("\n[CACH] - Initializing cache: %s (backend: %s)", name, c.Backend)
			redisMap[name] = &RedisCacheWrapper{
				Name:   name,
				Client: client,
				Ctx:    ctx,
			}

		//||------------------------------------------------------------------------------------------------||
		//|| KeyDB (redis protocol)
		//||------------------------------------------------------------------------------------------------||

		case "keydb", "KEYDB":
			client, ctx, err := connectRedis(c)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("cache '%s' keydb connect failed: %w", name, err)
			}
			keydbMap[name] = &RedisCacheWrapper{
				Name:   name,
				Client: client,
				Ctx:    ctx,
			}

		//||------------------------------------------------------------------------------------------------||
		//|| Memcached
		//||------------------------------------------------------------------------------------------------||

		case "memcached", "MEMCACHED":
			client, err := connectMemcached(c)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("cache '%s' memcached connect failed: %w", name, err)
			}
			memcMap[name] = &MemcachedCacheWrapper{
				Name:   name,
				Client: client,
			}

		//||------------------------------------------------------------------------------------------------||
		//|| In-Memory
		//||------------------------------------------------------------------------------------------------||

		case "memory", "MEMORY":
			memMap[name] = &MemoryCacheWrapper{
				Name:  name,
				Store: initMemoryStore(),
			}

		//||------------------------------------------------------------------------------------------------||
		//|| Unsupported
		//||------------------------------------------------------------------------------------------------||

		default:
			return nil, nil, nil, nil, fmt.Errorf("unsupported cache backend: %s", c.Backend)
		}
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Return Maps
	//||------------------------------------------------------------------------------------------------||

	return redisMap, keydbMap, memcMap, memMap, nil
}
