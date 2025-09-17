package db

import (
	"github.com/ralphferrara/aria/auth/setup"
	"gorm.io/gorm"
)

func AuthDB() *gorm.DB {
	if setup.Setup.Database == nil {
		panic("AuthDB: database connection not initialized")
	}
	return setup.Setup.Database
}
