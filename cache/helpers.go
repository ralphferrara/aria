//||------------------------------------------------------------------------------------------------||
//|| Cache Package: Helpers
//|| helpers.go
//||------------------------------------------------------------------------------------------------||

package cache

import (
	"context"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/ralphferrara/aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| Build Redis/KeyDB Options
//||------------------------------------------------------------------------------------------------||

func buildRedisOptions(cfg config.CacheInstanceConfig) *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Connect Redis/KeyDB
//||------------------------------------------------------------------------------------------------||

func connectRedis(cfg config.CacheInstanceConfig) (*redis.Client, context.Context, error) {
	opts := buildRedisOptions(cfg)
	ctx := context.Background()
	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Connect Memcached (Single or Multi-Server)
//||------------------------------------------------------------------------------------------------||

func connectMemcached(cfg config.CacheInstanceConfig) (*memcache.Client, error) {
	var servers []string
	if len(cfg.Servers) > 0 {
		servers = cfg.Servers
	} else {
		servers = []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}
	}
	client := memcache.New(servers...)
	err := client.Ping()
	if err != nil {
		return nil, err
	}
	return client, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Init Memory Store
//||------------------------------------------------------------------------------------------------||

func initMemoryStore() map[string]interface{} {
	return make(map[string]interface{})
}
