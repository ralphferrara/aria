//||------------------------------------------------------------------------------------------------||
//|| DB Package: Initialization
//|| init.go
//||------------------------------------------------------------------------------------------------||

package db

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"aria/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//||------------------------------------------------------------------------------------------------||
//|| DB: Globals (SQL & Mongo)
//||------------------------------------------------------------------------------------------------||

var (
	SQL   = map[string]*GormWrapper{}
	Mongo = map[string]*MongoWrapper{}
)

//||------------------------------------------------------------------------------------------------||
//|| DB: Init - Connects all DBs from config
//||------------------------------------------------------------------------------------------------||

func Init() error {
	cfg := config.GetConfig()
	for name, dbCfg := range cfg.DB {
		switch dbCfg.Driver {
		//||------------------------------------------------------------------------------------------------||
		//|| PostGres
		//||------------------------------------------------------------------------------------------------||
		case "postgres":
			dsn := buildDSN(dbCfg)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return fmt.Errorf("failed to connect to postgres '%s': %w", name, err)
			}
			SQL[name] = &GormWrapper{Name: name, DB: db}

		//||------------------------------------------------------------------------------------------------||
		//|| MySQL, MariaDB
		//||------------------------------------------------------------------------------------------------||
		case "mysql", "mariadb":
			dsn := buildDSN(dbCfg)
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				return fmt.Errorf("failed to connect to mysql/mariadb '%s': %w", name, err)
			}
			SQL[name] = &GormWrapper{Name: name, DB: db}

		//||------------------------------------------------------------------------------------------------||
		//|| MongoDB
		//||------------------------------------------------------------------------------------------------||
		case "mongo":
			mCfg := config.GetConfig().DB[name]
			mdb, err := connectMongo(mCfg)
			if err != nil {
				return fmt.Errorf("mongo connect failed for '%s': %w", name, err)
			}
			Mongo[name] = &MongoWrapper{Name: name, Database: mdb}

		//||------------------------------------------------------------------------------------------------||
		//|| Default
		//||------------------------------------------------------------------------------------------------||
		default:
			return fmt.Errorf("unsupported db driver: %s", dbCfg.Driver)
		}
	}
	return nil
}
