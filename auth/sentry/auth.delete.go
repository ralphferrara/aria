package sentry

import (
	"base/db/abstract"
	"net/http"

	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/auth/actions"
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
	expected := actions.GenerateQuitToken(account.PrivateHash)
	//||------------------------------------------------------------------------------------------------||
	//|| Quit Token
	//||------------------------------------------------------------------------------------------------||
	if clientToken != expected {
		responses.Error(w, http.StatusForbidden, "Invalid token")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Quit Token
	//||------------------------------------------------------------------------------------------------||
	if err := abstract.DeleteAccount(account.ID); err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to delete account")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Success
	//||------------------------------------------------------------------------------------------------||
	responses.Success(w, http.StatusOK, responseDelete{
		Message: "Account deleted successfully",
	})
}
