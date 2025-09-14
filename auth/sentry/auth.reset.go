package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"net/http"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/base/validate"
	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/auth/db"
)

// ||------------------------------------------------------------------------------------------------||
// || ResetPasswordHandler
// || Allows a logged-in user to reset their password by providing and confirming a new password
// ||------------------------------------------------------------------------------------------------||

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| Get Session Cookie
	//||------------------------------------------------------------------------------------------------||

	cookie, err := r.Cookie("session")
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "No session cookie")
		return
	}

	session, err := actions.FetchSession(cookie.Value)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Get Account
	//||------------------------------------------------------------------------------------------------||

	account, err := db.GetAccountByID(fmt.Sprintf("%d", session.ID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to fetch account")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Parse Password Fields
	//||------------------------------------------------------------------------------------------------||

	password := r.FormValue("password")
	confirm := r.FormValue("confirmPassword")

	error := validate.IsValidPassword(password)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error.Error())
		return
	}

	if password != confirm {
		responses.Error(w, http.StatusBadRequest, app.Err("Auth").Code("PASSWORD_MISMATCH"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Hash Password and Generate Salt
	//||------------------------------------------------------------------------------------------------||

	passwordHash := actions.GeneratePasswordWithSalt(password, account.Salt)
	if passwordHash == "" {
		responses.Error(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Update Account Record
	//||------------------------------------------------------------------------------------------------||

	err = db.UpdatePassword(account.ID, passwordHash)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Success
	//||------------------------------------------------------------------------------------------------||
	responses.Success(w, http.StatusOK, map[string]any{
		"message": "Password successfully updated",
	})
}
