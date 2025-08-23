//||------------------------------------------------------------------------------------------------||
//|| DB Package: Helpers
//|| helpers.go
//||------------------------------------------------------------------------------------------------||

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ralphferrara/aria/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//||------------------------------------------------------------------------------------------------||
//|| Build DSN
//||------------------------------------------------------------------------------------------------||

func buildDSN(cfg config.DBInstanceConfig) string {
	switch cfg.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
	case "mysql", "mariadb":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	}
	return ""
}

//||------------------------------------------------------------------------------------------------||
//|| Build DSN
//||------------------------------------------------------------------------------------------------||

func connectMongo(cfg config.DBInstanceConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(cfg.Host)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client.Database(cfg.Database), nil
}
