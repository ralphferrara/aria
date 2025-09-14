package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"net/http"

	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/auth/setup"
)

//||------------------------------------------------------------------------------------------------||
//|| Handler :: Initiates Forgot Password Flow
//||------------------------------------------------------------------------------------------------||

func AuthMeHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| Get the Session Cookie
	//||------------------------------------------------------------------------------------------------||

	_, account, session, err := actions.LoadSessionAccount(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Check Session
	//||------------------------------------------------------------------------------------------------||

	if account.ID <= 0 || session.ID != account.ID {
		responses.Error(w, http.StatusUnauthorized, app.Err("Auth").Code("INVALID_SESSION"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Level
	//||------------------------------------------------------------------------------------------------||

	if account.Level < 1 {
		responses.Error(w, http.StatusForbidden, app.Err("Auth").Code("INSUFFICIENT_LEVEL"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Banned
	//||------------------------------------------------------------------------------------------------||

	if account.Status == app.Constants("AccountStatus").Code("Suspended") {
		responses.Error(w, http.StatusForbidden, app.Err("Auth").Code("ACCOUNT_SUSPENDED"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Deleted
	//||------------------------------------------------------------------------------------------------||

	if account.Status == app.Constants("AccountStatus").Code("Deleted") {
		responses.Error(w, http.StatusForbidden, app.Err("Auth").Code("ACCOUNT_DELETED"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Success Response
	//||------------------------------------------------------------------------------------------------||

	setup.Setup.Functions.OnAuthCheck(w, r, account.ID, session.Identifier)

}
