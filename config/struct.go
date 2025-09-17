//||------------------------------------------------------------------------------------------------||
//|| Config Package: Types
//|| types.go
//||------------------------------------------------------------------------------------------------||

package config

//||------------------------------------------------------------------------------------------------||
//|| Config: Root Struct
//||------------------------------------------------------------------------------------------------||

type Config struct {
	App      AppConfig                        `json:"app"`
	DB       map[string]DBInstanceConfig      `json:"db"`
	Cache    map[string]CacheInstanceConfig   `json:"cache"`
	Storage  map[string]StorageInstanceConfig `json:"storage"`
	Queue    map[string]QueueInstanceConfig   `json:"queue"`
	HTTP     map[string]HTTPInstanceConfig    `json:"http"` // <-- add this
	Auth     AuthConfig                       `json:"auth"`
	Locale   LocaleConfig                     `json:"locale"`
	Template TemplateConfig                   `json:"template"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: App Section
//||------------------------------------------------------------------------------------------------||

type AppConfig struct {
	Name  string `json:"name"`
	Env   string `json:"env"`
	Debug bool   `json:"debug"`
	Port  int    `json:"port"`
	Salt  string `json:"salt"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: DB Section
//||------------------------------------------------------------------------------------------------||

type DBInstanceConfig struct {
	Driver   string `json:"driver"` // postgres | mysql | mariadb | mongo
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	SSLMode  string `json:"sslmode,omitempty"` // postgres only
	URI      string `json:"uri,omitempty"`     // optional mongo URI
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Cache Section
//||------------------------------------------------------------------------------------------------||

type CacheInstanceConfig struct {
	Backend  string   `json:"backend"` // redis | keydb | memcached | memory
	Host     string   `json:"host,omitempty"`
	Port     int      `json:"port,omitempty"`
	Password string   `json:"password,omitempty"`
	DB       int      `json:"db,omitempty"`
	Servers  []string `json:"servers,omitempty"` // memcached
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Storage Section
//||------------------------------------------------------------------------------------------------||

type StorageInstanceConfig struct {
	Backend   string `json:"backend"` // s3 | minio | gcs | local
	Bucket    string `json:"bucket,omitempty"`
	Region    string `json:"region,omitempty"`
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	Dir       string `json:"dir,omitempty"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Queue Section
//||------------------------------------------------------------------------------------------------||

type QueueInstanceConfig struct {
	Backend  string `json:"backend"` // rabbitmq, etc.
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Vhost    string `json:"vhost,omitempty"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: HTTP Section
//||------------------------------------------------------------------------------------------------||

type HTTPInstanceConfig struct {
	Backend      string `json:"backend"` // "mux" or "http"
	Port         int    `json:"port"`
	Cors         bool   `json:"cors"`
	Middleware   bool   `json:"middleware"`
	ErrorHandler bool   `json:"error_handler"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Locale Section
//||------------------------------------------------------------------------------------------------||

type LocaleConfig struct {
	Default   string   `json:"default"`
	Supported []string `json:"supported"`
	Directory string   `json:"directory"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Template Section
//||------------------------------------------------------------------------------------------------||

type TemplateConfig struct {
	Dir   string `json:"dir"`
	Cache bool   `json:"cache"`
}
