//||------------------------------------------------------------------------------------------------||
//|| DB Package: Initialization
//|| init.go
//||------------------------------------------------------------------------------------------------||

package db

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"

	"github.com/ralphferrara/aria/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//||------------------------------------------------------------------------------------------------||
//|| DB: Init - Connects all DBs from config
//||------------------------------------------------------------------------------------------------||

func Init(cfg *config.Config) (map[string]*GormWrapper, map[string]*MongoWrapper, error) {

	sqlDB := make(map[string]*GormWrapper)
	mongoDB := make(map[string]*MongoWrapper)

	for name, dbCfg := range cfg.DB {
		switch dbCfg.Driver {

		//||------------------------------------------------------------------------------------------------||
		//|| PostGres
		//||------------------------------------------------------------------------------------------------||
		case "postgres":
			dsn := buildDSN(dbCfg)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, nil, fmt.Errorf("failed to connect to postgres '%s': %w", name, err)
			}
			sqlDB[name] = &GormWrapper{Name: name, DB: db}
			fmt.Printf("\n[ DB ] - Initialized database: %s (backend: %s)", name, dbCfg.Driver)

		//||------------------------------------------------------------------------------------------------||
		//|| MySQL, MariaDB
		//||------------------------------------------------------------------------------------------------||
		case "mysql", "mariadb":
			dsn := buildDSN(dbCfg)
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, nil, fmt.Errorf("failed to connect to mysql/mariadb '%s': %w", name, err)
			}
			sqlDB[name] = &GormWrapper{Name: name, DB: db}
			fmt.Printf("\n[ DB ] - Initialized database: %s (backend: %s)", name, dbCfg.Driver)
		//||------------------------------------------------------------------------------------------------||
		//|| MongoDB
		//||------------------------------------------------------------------------------------------------||
		case "mongo":
			mdb, err := connectMongo(dbCfg)
			if err != nil {
				return nil, nil, fmt.Errorf("mongo connect failed for '%s': %w", name, err)
			}
			mongoDB[name] = &MongoWrapper{Name: name, Database: mdb}
			fmt.Printf("\n[ DB ] - Initialized database: %s (backend: %s)", name, dbCfg.Driver)

		//||------------------------------------------------------------------------------------------------||
		//|| Default
		//||------------------------------------------------------------------------------------------------||
		default:
			return nil, nil, fmt.Errorf("unsupported db driver: %s", dbCfg.Driver)
		}
	}

	return sqlDB, mongoDB, nil
}
