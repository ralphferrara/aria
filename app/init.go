//||------------------------------------------------------------------------------------------------||
//|| App Package: Loader & Bootstrapper
//|| init.go
//||------------------------------------------------------------------------------------------------||

package app

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"github.com/ralphferrara/aria/cache"
	"github.com/ralphferrara/aria/config"
	"github.com/ralphferrara/aria/db"
	"github.com/ralphferrara/aria/log"
	"github.com/ralphferrara/aria/queue"
	"github.com/ralphferrara/aria/storage"
)

//||------------------------------------------------------------------------------------------------||
//|| App: Globals
//||------------------------------------------------------------------------------------------------||

var (
	Config         *config.Config
	Storages       = map[string]*storage.Storage{}
	SQLDB          = map[string]*db.GormWrapper{}
	MongoDB        = map[string]*db.MongoWrapper{}
	QueueRabbit    = map[string]*queue.RabbitMQWrapper{}
	CacheRedis     = map[string]*cache.RedisCacheWrapper{}
	CacheKeyDB     = map[string]*cache.RedisCacheWrapper{}
	CacheMemcached = map[string]*cache.MemcachedCacheWrapper{}
	CacheMemory    = map[string]*cache.MemoryCacheWrapper{}
	Log            = log.Log
)

//||------------------------------------------------------------------------------------------------||
//|| App: Init
//||------------------------------------------------------------------------------------------------||

func Init(configFile string) error {
	//||------------------------------------------------------------------------------------------------||
	//|| Load Config
	//||------------------------------------------------------------------------------------------------||
	cfg, err := config.Init(configFile)
	if err != nil {
		Log.Error("app", "Failed to load config: %v", err)
		return err
	}
	Log.Init(cfg)
	Log.Info("app", "Config loaded from %s", configFile)
	Config = cfg

	//||------------------------------------------------------------------------------------------------||
	//|| Init Storage(s)
	//||------------------------------------------------------------------------------------------------||
	for name, sCfg := range cfg.Storage {
		storeCfg := storage.ConvertFromConfig(sCfg)
		st := &storage.Storage{Config: storeCfg}
		if err := st.Init(); err != nil {
			Log.Error("app", "Failed to init storage '%s': %v", name, err)
			return err
		}
		Storages[name] = st
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Init Databases
	//||------------------------------------------------------------------------------------------------||
	if err := db.Init(); err != nil {
		Log.Error("app", "Failed to init databases: %v", err)
		return err
	}
	SQLDB = db.SQL
	MongoDB = db.Mongo

	//||------------------------------------------------------------------------------------------------||
	//|| Init Queues
	//||------------------------------------------------------------------------------------------------||
	if err := queue.Init(); err != nil {
		Log.Error("app", "Failed to init queues: %v", err)
		return err
	}
	QueueRabbit = queue.Rabbit

	//||------------------------------------------------------------------------------------------------||
	//|| Init Caches
	//||------------------------------------------------------------------------------------------------||
	if err := cache.Init(); err != nil {
		Log.Error("app", "Failed to init caches: %v", err)
		return err
	}
	CacheRedis = cache.Redis
	CacheKeyDB = cache.KeyDB
	CacheMemcached = cache.Memcached
	CacheMemory = cache.Memory

	Log.Info("app", "Aria app initialized successfully")
	return nil
}
