//||------------------------------------------------------------------------------------------------||
//|| App Package: Loader & Bootstrapper
//|| init.go
//||------------------------------------------------------------------------------------------------||

package app

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ralphferrara/aria/cache"
	"github.com/ralphferrara/aria/config"
	"github.com/ralphferrara/aria/db"
	"github.com/ralphferrara/aria/http"
	"github.com/ralphferrara/aria/locale"
	"github.com/ralphferrara/aria/log"
	"github.com/ralphferrara/aria/queue"
	"github.com/ralphferrara/aria/storage"
)

//||------------------------------------------------------------------------------------------------||
//|| App: Globals
//||------------------------------------------------------------------------------------------------||

var (
	Config         *config.Config
	HTTP           map[string]*http.HTTPWrapper
	Storages       map[string]*storage.Storage
	SQLDB          map[string]*db.GormWrapper
	MongoDB        map[string]*db.MongoWrapper
	QueueRabbit    map[string]*queue.RabbitMQWrapper
	CacheRedis     map[string]*cache.RedisCacheWrapper
	CacheKeyDB     map[string]*cache.RedisCacheWrapper
	CacheMemcached map[string]*cache.MemcachedCacheWrapper
	CacheMemory    map[string]*cache.MemoryCacheWrapper
	Log            log.Logger
	Locales        locale.LocaleWrapper
)

//||------------------------------------------------------------------------------------------------||
//|| App: Init (package-level, no AriaApplication wrapper)
//||------------------------------------------------------------------------------------------------||

func Init(configFile string) {

	//||------------------------------------------------------------------------------------------------||
	//|| Load a Version
	//||------------------------------------------------------------------------------------------------||

	fmt.Print("\n\n")
	fmt.Print("\033[44;97m") // Blue background, white text
	fmt.Println("||------------------------------------------------------||")
	fmt.Println("|| Welcome to Aria")
	fmt.Println("||------------------------------------------------------||\033[0m")
	fmt.Println("") // Reset

	//||------------------------------------------------------------------------------------------------||
	//|| Logger
	//||------------------------------------------------------------------------------------------------||

	Log = log.Init("aria")

	//||------------------------------------------------------------------------------------------------||
	//|| Load Config
	//||------------------------------------------------------------------------------------------------||

	cfg, err := config.Init(configFile)
	if err != nil {
		Log.Error("app", "Failed to load config: %v", err)
		os.Exit(1)
	}
	Config = cfg
	fmt.Printf("[CNFG] - Config loaded from %s", configFile)

	//||------------------------------------------------------------------------------------------------||
	//|| Constants (package-level)
	//||------------------------------------------------------------------------------------------------||

	InitConstants()

	//||------------------------------------------------------------------------------------------------||
	//|| HTTP(s)
	//||------------------------------------------------------------------------------------------------||

	httpMap, err := http.Init(cfg)
	if err != nil {
		Log.Error("app", "Failed to init HTTP server(s): %v", err)
		os.Exit(1)
	}
	HTTP = httpMap

	//||------------------------------------------------------------------------------------------------||
	//|| Storages
	//||------------------------------------------------------------------------------------------------||

	stMap, err := storage.Init(cfg)
	if err != nil {
		Log.Error("app", "Failed to init storage(s): %v", err)
		os.Exit(1)
	}
	Storages = stMap

	//||------------------------------------------------------------------------------------------------||
	//|| Databases
	//||------------------------------------------------------------------------------------------------||

	sqlMap, mongoMap, err := db.Init(cfg)
	if err != nil {
		Log.Error("app", "Failed to init database(s): %v", err)
		os.Exit(1)
	}
	SQLDB = sqlMap
	MongoDB = mongoMap

	//||------------------------------------------------------------------------------------------------||
	//|| Queues
	//||------------------------------------------------------------------------------------------------||

	qMap, err := queue.Init(cfg)
	if err != nil {
		Log.Error("\nFailed to init queue(s): %v", err)
		os.Exit(1)
	}
	QueueRabbit = qMap

	//||------------------------------------------------------------------------------------------------||
	//|| Locales
	//||------------------------------------------------------------------------------------------------||

	locale.Init(Config.Locale.Directory)

	//||------------------------------------------------------------------------------------------------||
	//|| Caches
	//||------------------------------------------------------------------------------------------------||

	rMap, kdbMap, mcMap, memMap, err := cache.Init(cfg)
	if err != nil {
		Log.Error("app", "Failed to init cache(s): %v", err)
		os.Exit(1)
	}
	CacheRedis = rMap
	CacheKeyDB = kdbMap
	CacheMemcached = mcMap
	CacheMemory = memMap
}

//||------------------------------------------------------------------------------------------------||
//|| Listen
//||------------------------------------------------------------------------------------------------||

func ListenForShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		Log.Info("Received signal: %s", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		Shutdown(ctx)
		os.Exit(0)
	}()
}

//||------------------------------------------------------------------------------------------------||
//|| App: Shutdown (graceful best-effort) - package-level
//||------------------------------------------------------------------------------------------------||

func Shutdown(ctx context.Context) {

	//||------------------------------------------------------------------------------------------------||
	//|| Setup Deadline
	//||------------------------------------------------------------------------------------------------||

	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		dctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ctx = dctx
		deadline, _ = ctx.Deadline()
	}

	Log.Info("Shutdown initiated (deadline: %s)", deadline.Format(time.RFC3339))

	//||------------------------------------------------------------------------------------------------||
	//|| Caches
	//||------------------------------------------------------------------------------------------------||

	for name, c := range CacheRedis {
		closeIf(c)
		Log.Info("CacheRedis '%s' closed", name)
	}
	for name, c := range CacheKeyDB {
		closeIf(c)
		Log.Info("CacheKeyDB '%s' closed", name)
	}
	for name, c := range CacheMemcached {
		closeIf(c)
		Log.Info("CacheMemcached '%s' closed", name)
	}
	for name, c := range CacheMemory {
		closeIf(c)
		Log.Info("CacheMemory '%s' closed", name)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Queues
	//||------------------------------------------------------------------------------------------------||

	for name, q := range QueueRabbit {
		closeIf(q)
		Log.Info("QueueRabbit '%s' closed", name)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Databases
	//||------------------------------------------------------------------------------------------------||

	for name, g := range SQLDB {
		closeIf(g)
		Log.Info("SQLDB '%s' closed", name)
	}
	for name, m := range MongoDB {
		closeIf(m)
		Log.Info("MongoDB '%s' closed", name)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Storages
	//||------------------------------------------------------------------------------------------------||

	for name, s := range Storages {
		closeIf(s)
		Log.Info("Storage '%s' closed", name)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| HTTP servers (graceful)
	//||------------------------------------------------------------------------------------------------||

	for name, h := range HTTP {
		if h.Server != nil {
			_ = h.Server.Shutdown(ctx)
		} else {
			closeIf(h)
		}
		Log.Info("HTTP '%s' closed", name)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Done
	//||------------------------------------------------------------------------------------------------||

	Log.Info("Shutdown complete")
}

//||------------------------------------------------------------------------------------------------||
//|| helper: closeIf
//||------------------------------------------------------------------------------------------------||

type closer interface{ Close() error }
type stopper interface{ Stop() error }

func closeIf(x any) {
	switch v := x.(type) {
	case nil:
		return
	case closer:
		_ = v.Close()
	case stopper:
		_ = v.Stop()
	}
}
