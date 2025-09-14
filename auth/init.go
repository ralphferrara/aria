package auth

import (
	"github.com/ralphferrara/aria/auth/setup"
	"github.com/ralphferrara/aria/config"
	"gorm.io/gorm"
)

//||------------------------------------------------------------------------------------------------||
//|| Init
//||------------------------------------------------------------------------------------------------||

func Init(gormDB *gorm.DB, config config.AuthConfig) {
	setup.Setup.Initialized = true
	setup.Setup.CSRF = config.CSRF
	setup.Setup.Pepper = config.Pepper
	setup.Setup.Database = gormDB
	setup.Setup.Table = config.Table
	InitConstants()
}
