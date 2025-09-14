package sentry

import (
	"net/http"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/auth/db"
	"github.com/ralphferrara/aria/auth/setup"
)

//||------------------------------------------------------------------------------------------------||
//|| Response
//||------------------------------------------------------------------------------------------------||

type responseDelete struct {
	Message string `json:"message"`
}

//||------------------------------------------------------------------------------------------------||
//|| DeleteAccountHandler – validates token and deletes account
//||------------------------------------------------------------------------------------------------||

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	//||------------------------------------------------------------------------------------------------||
	//|| Get the Session Cookie
	//||------------------------------------------------------------------------------------------------||
	clientToken := r.FormValue("quitToken")
	_, account, _, err := actions.LoadSessionAccount(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Quit Token
	//||------------------------------------------------------------------------------------------------||
	expected := actions.GenerateQuitToken(account.Password)
	//||------------------------------------------------------------------------------------------------||
	//|| Quit Token
	//||------------------------------------------------------------------------------------------------||
	if clientToken != expected {
		responses.Error(w, http.StatusForbidden, app.Err("Auth").Code("ACCOUNT_TOKEN_MISMATCH"))
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Call External Before Delete Hook
	//||------------------------------------------------------------------------------------------------||
	before := setup.Setup.Functions.OnBeforeAccountDelete(account.ID)
	if before != nil {
		responses.Error(w, http.StatusInternalServerError, before.Error())
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Call Account Delete
	//||------------------------------------------------------------------------------------------------||
	if err := db.DeleteAccount(account.ID); err != nil {
		responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("ACCOUNT_DELETE_FAILED"))
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Call External Before Delete Hook
	//||------------------------------------------------------------------------------------------------||
	after := setup.Setup.Functions.OnAfterAccountDelete(account.ID)
	if after != nil {
		responses.Error(w, http.StatusInternalServerError, after.Error())
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Success
	//||------------------------------------------------------------------------------------------------||
	responses.Success(w, http.StatusOK, responseDelete{
		Message: "OK",
	})
}
