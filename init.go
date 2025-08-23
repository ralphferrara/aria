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
	Log      = log.Print
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
		log.Error("app", "Failed to load config: %v", err)
		return err
	}
	log.Init(cfg)
	log.Info("app", "Config loaded from %s", configFile)
	Config = cfg
	//||------------------------------------------------------------------------------------------------||
	//|| Init Storage(s)
	//||------------------------------------------------------------------------------------------------||

	for name, sCfg := range cfg.Storage {
		storeCfg := storage.ConvertFromConfig(sCfg)
		st := &storage.Storage{Config: storeCfg}
		if err := st.Init(); err != nil {
			log.Error("app", "Failed to init storage '%s': %v", name, err)
			return err
		}
		Storages[name] = st
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Init Databases
	//||------------------------------------------------------------------------------------------------||
	if err := db.Init(); err != nil {
		log.Error("app", "Failed to init databases: %v", err)
		return err
	}
	// Optionally assign/copy for app-global shortcut
	SQLDB = db.SQL     // SQLDB["main"], etc.
	MongoDB = db.Mongo // MongoDB["mongo"], etc.

	log.Info("app", "Aria app initialized successfully")
	return nil
}
