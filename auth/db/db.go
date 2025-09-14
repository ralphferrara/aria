package db

import (
	"github.com/ralphferrara/aria/auth/setup"
	"gorm.io/gorm"
)

func AuthDB() *gorm.DB {
	return setup.Setup.Database
}
