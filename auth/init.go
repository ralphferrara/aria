package auth

import (
	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/setup/authentication"
)

//||------------------------------------------------------------------------------------------------||
//|| Setup
//||------------------------------------------------------------------------------------------------||

var (
	Setup = authentication.Setup{
		Initialized: false,
		Pepper:      "",
		CSRF:        "",
		Database:    "",
	}
)

//||------------------------------------------------------------------------------------------------||
//|| Init
//||------------------------------------------------------------------------------------------------||

func Init() {
	Setup.Initialized = true
	Setup.CSRF = app.Config.Auth.CSRF
	Setup.Pepper = app.Config.Auth.Pepper
	Setup.Database = app.Config.Auth.Database
}
