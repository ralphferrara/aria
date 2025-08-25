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
}

//||------------------------------------------------------------------------------------------------||
//|| Config: DB Section
//||------------------------------------------------------------------------------------------------||

type DBInstanceConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode,omitempty"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Cache Section
//||------------------------------------------------------------------------------------------------||

type CacheInstanceConfig struct {
	Backend  string   `json:"backend"`
	Host     string   `json:"host,omitempty"`
	Port     int      `json:"port,omitempty"`
	Password string   `json:"password,omitempty"`
	DB       int      `json:"db,omitempty"`
	Servers  []string `json:"servers,omitempty"` // <-- add this
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Storage Section
//||------------------------------------------------------------------------------------------------||

type StorageInstanceConfig struct {
	Backend   string `json:"backend"`
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
	Backend  string `json:"backend"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Vhost    string `json:"vhost,omitempty"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Auth Section
//||------------------------------------------------------------------------------------------------||

type AuthConfig struct {
	JwtSecret     string `json:"jwt_secret"`
	SessionExpiry int    `json:"session_expiry"`
	TokenIssuer   string `json:"token_issuer"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Locale Section
//||------------------------------------------------------------------------------------------------||

type LocaleConfig struct {
	Default   string   `json:"default"`
	Supported []string `json:"supported"`
}

//||------------------------------------------------------------------------------------------------||
//|| Config: Template Section
//||------------------------------------------------------------------------------------------------||

type TemplateConfig struct {
	Dir   string `json:"dir"`
	Cache bool   `json:"cache"`
}
