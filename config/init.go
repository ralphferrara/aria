//||------------------------------------------------------------------------------------------------||
//|| Config Package: Loader
//|| init.go
//||------------------------------------------------------------------------------------------------||

package config

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

//||------------------------------------------------------------------------------------------------||
//|| Config: Globals
//||------------------------------------------------------------------------------------------------||

var (
	cfg  *Config
	once sync.Once
	mu   sync.RWMutex
)

//||------------------------------------------------------------------------------------------------||
//|| Init
//||------------------------------------------------------------------------------------------------||

func Init(
	path string,
) (*Config, error) {

	var loadErr error

	once.Do(func() {

		//||------------------------------------------------------------------------------------------------||
		//|| Load JSON Config
		//||------------------------------------------------------------------------------------------------||

		file, err := os.Open(path)
		if err != nil {
			loadErr = fmt.Errorf("failed to open config file: %w", err)
			return
		}
		defer file.Close()

		local := &Config{}
		dec := json.NewDecoder(file)
		dec.DisallowUnknownFields()

		if err := dec.Decode(local); err != nil {
			loadErr = fmt.Errorf("decode config: %w", err)
			return
		}

		//||------------------------------------------------------------------------------------------------||
		//|| Expand Env Vars
		//||------------------------------------------------------------------------------------------------||

		expandEnvStrings(local)

		//||------------------------------------------------------------------------------------------------||
		//|| Print Debug Config Dump (optional)
		//||------------------------------------------------------------------------------------------------||

		//PrintConsoleConfig(local)

		//||------------------------------------------------------------------------------------------------||
		//|| Normalize + Validate
		//||------------------------------------------------------------------------------------------------||

		normalize(local)
		if err := validate(local); err != nil {
			loadErr = fmt.Errorf("invalid config: %w", err)
			return
		}

		//||------------------------------------------------------------------------------------------------||
		//|| Assign to Global Singleton
		//||------------------------------------------------------------------------------------------------||

		mu.Lock()
		cfg = local
		mu.Unlock()
	})

	//||------------------------------------------------------------------------------------------------||
	//|| Final Check & Return
	//||------------------------------------------------------------------------------------------------||

	if loadErr != nil {
		return nil, loadErr
	}

	mu.RLock()
	defer mu.RUnlock()

	if cfg == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	return cfg, nil
}

//||------------------------------------------------------------------------------------------------||
//|| GetConfig: Accessor (returns pointer, may be nil if not loaded)
//||------------------------------------------------------------------------------------------------||

func InProduction() bool {
	return (cfg.App.Env == "production")
}

//||------------------------------------------------------------------------------------------------||
//|| GetConfig: Accessor (returns pointer, may be nil if not loaded)
//||------------------------------------------------------------------------------------------------||

func GetConfig() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return cfg
}

//||------------------------------------------------------------------------------------------------||
//|| Must: Convenience accessor (panics if config is not initialized)
//||------------------------------------------------------------------------------------------------||

func Must() *Config {
	c := GetConfig()
	if c == nil {
		panic("config not initialized: call config.Init(path) first")
	}
	return c
}

//||------------------------------------------------------------------------------------------------||
//|| Reset: TEST-ONLY helper to clear the singleton (not for production)
//||------------------------------------------------------------------------------------------------||

func Reset() {
	mu.Lock()
	defer mu.Unlock()
	cfg = nil
	// allow re-Init after Reset
	once = sync.Once{}
}

//||------------------------------------------------------------------------------------------------||
//|| expandEnvStrings walks string fields and applies os.ExpandEnv
//||------------------------------------------------------------------------------------------------||

func expandEnvStrings(c *Config) {
	// App
	c.App.Name = os.ExpandEnv(c.App.Name)
	c.App.Env = os.ExpandEnv(c.App.Env)

	// DB
	for k, v := range c.DB {
		v.Driver = os.ExpandEnv(v.Driver)
		v.Host = os.ExpandEnv(v.Host)
		v.User = os.ExpandEnv(v.User)
		v.Password = os.ExpandEnv(v.Password)
		v.Database = os.ExpandEnv(v.Database)
		v.SSLMode = os.ExpandEnv(v.SSLMode)
		c.DB[k] = v
	}

	// Cache
	for k, v := range c.Cache {
		v.Backend = os.ExpandEnv(v.Backend)
		v.Host = os.ExpandEnv(v.Host)
		v.Password = os.ExpandEnv(v.Password)
		c.Cache[k] = v
	}

	// Storage
	for k, v := range c.Storage {
		v.Backend = os.ExpandEnv(v.Backend)
		v.Bucket = os.ExpandEnv(v.Bucket)
		v.Region = os.ExpandEnv(v.Region)
		v.AccessKey = os.ExpandEnv(v.AccessKey)
		v.SecretKey = os.ExpandEnv(v.SecretKey)
		v.Endpoint = os.ExpandEnv(v.Endpoint)
		v.Dir = os.ExpandEnv(v.Dir)
		c.Storage[k] = v
	}

	// Queue
	for k, v := range c.Queue {
		v.Backend = os.ExpandEnv(v.Backend)
		v.Host = os.ExpandEnv(v.Host)
		v.User = os.ExpandEnv(v.User)
		v.Password = os.ExpandEnv(v.Password)
		v.Vhost = os.ExpandEnv(v.Vhost)
		c.Queue[k] = v
	}

	// HTTP
	for k, v := range c.HTTP {
		v.Backend = os.ExpandEnv(v.Backend)
		c.HTTP[k] = v
	}

	// Auth
	c.Auth.CSRF = os.ExpandEnv(c.Auth.CSRF)
	c.Auth.Pepper = os.ExpandEnv(c.Auth.Pepper)
	c.Auth.Table = os.ExpandEnv(c.Auth.Table)

	// Locale
	c.Locale.Default = os.ExpandEnv(c.Locale.Default)

	// Template
	c.Template.Dir = os.ExpandEnv(c.Template.Dir)
}

//||------------------------------------------------------------------------------------------------||
//|| normalize: tidy up enums/strings (e.g., backend types)
//||------------------------------------------------------------------------------------------------||

func normalize(c *Config) {
	// http backend normalize
	for k, v := range c.HTTP {
		switch strings.ToLower(v.Backend) {
		case "mux":
			v.Backend = "mux"
		case "http", "servemux", "nethttp":
			v.Backend = "http"
		default:
			// keep as-is; validate() will flag if unsupported
		}
		c.HTTP[k] = v
	}
}
