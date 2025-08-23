//||------------------------------------------------------------------------------------------------||
//|| App Package: Loader & Bootstrapper
//|| init.go
//||------------------------------------------------------------------------------------------------||

package aria

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"aria/config"
	"aria/db"
	"aria/log"
	"aria/storage"
)

//||------------------------------------------------------------------------------------------------||
//|| App: Globals
//||------------------------------------------------------------------------------------------------||

var (
	Config   *config.Config
	Storages = map[string]*storage.Storage{}
	SQLDB    = map[string]*db.GormWrapper{} // <-- FIXED HERE
	MongoDB  = map[string]*db.MongoWrapper{}
	Log      = log.Log
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
	SQLDB = db.SQL     // SQLDB["main"], etc.
	MongoDB = db.Mongo // MongoDB["mongo"], etc.

	Log.Info("app", "Aria app initialized successfully")
	return nil
}
