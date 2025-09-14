package actions

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"net/http"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/db"
	"github.com/ralphferrara/aria/auth/types"
)

//||------------------------------------------------------------------------------------------------||
//|| Load Session Account
//||------------------------------------------------------------------------------------------------||

func LoadSessionAccount(r *http.Request) (http.Cookie, db.ModelAccount, types.SessionRecord, error) {

	//||------------------------------------------------------------------------------------------------||
	//|| Get the Session Cookie
	//||------------------------------------------------------------------------------------------------||

	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		return http.Cookie{}, db.ModelAccount{}, types.SessionRecord{}, app.Err("Auth").Error("MISSING_SESSION_COOKIE")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Get Session
	//||------------------------------------------------------------------------------------------------||

	session, err := FetchSession(cookie.Value)
	if err != nil {
		return *cookie, db.ModelAccount{}, types.SessionRecord{}, app.Err("Auth").Error("SESSION_LOOKUP_FAILED")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Get Database Account
	//||------------------------------------------------------------------------------------------------||

	account, err := db.GetAccountByID(fmt.Sprintf("%d", session.ID))
	if err != nil || account == nil {
		return *cookie, db.ModelAccount{}, types.SessionRecord{}, app.Err("Auth").Error("ACCOUNT_LOOKUP_FAILED")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Done
	//||------------------------------------------------------------------------------------------------||

	return *cookie, *account, session, nil
}
