package setup

import (
	"net/http"

	"github.com/ralphferrara/aria/auth/types"
	"gorm.io/gorm"
)

//||------------------------------------------------------------------------------------------------||
//|| Auth Setup
//||------------------------------------------------------------------------------------------------||

var Setup AuthSetup

//||------------------------------------------------------------------------------------------------||
//|| Auth: Globals
//||------------------------------------------------------------------------------------------------||

type AuthSetup struct {
	Initialized bool
	Domain      string
	Pepper      string
	CSRF        string
	Database    *gorm.DB
	Table       string
	Functions   AuthFunctions
}

//||------------------------------------------------------------------------------------------------||
//|| Auth Functions
//||------------------------------------------------------------------------------------------------||

type AuthFunctions struct {
	OnAccountComplete     func(r *http.Request, accountID int64, accountIdentifier string) error
	OnAuthCheck           func(w http.ResponseWriter, r *http.Request, authMe types.AuthMeRecord) error
	OnBeforeAccountDelete func(accountID int64) error
	OnAfterAccountDelete  func(accountID int64) error
}
