package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"net/http"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/auth/db"
	"github.com/ralphferrara/aria/base/validate"
	"github.com/ralphferrara/aria/responses"
)

//||------------------------------------------------------------------------------------------------||
//|| Response
//||------------------------------------------------------------------------------------------------||

type authCompleteResponse struct {
	Message string `json:"message"`
	Next    string `json:"next"`
}

//||------------------------------------------------------------------------------------------------||
//|| Handler
//||------------------------------------------------------------------------------------------------||

func CompleteHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| DB Account
	//||------------------------------------------------------------------------------------------------||

	cookie, dbAccount, session, err := actions.LoadSessionAccount(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Check Account Status
	//||------------------------------------------------------------------------------------------------||

	if dbAccount.Status != app.Constants("AccountStatus").Code("Pending") {
		responses.Error(w, http.StatusForbidden, app.Err("Auth").Code("ACCOUNT_ALREADY_CREATED"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||

	password := r.FormValue("password")

	//||------------------------------------------------------------------------------------------------||
	//||
	//|| Sanitize and Validate
	//|| Also generate the private/public key if needed
	//||
	//||------------------------------------------------------------------------------------------------||

	vp := validate.IsValidPassword(password)
	if vp != nil {
		responses.Error(w, http.StatusBadRequest, vp.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Password/Salt
	//||------------------------------------------------------------------------------------------------||

	passwordHash, saltHash := actions.GeneratePassword(password)
	if passwordHash == "" {
		responses.Error(w, http.StatusBadRequest, app.Err("Auth").Code("PASSWORD_GEN_FAILED"))
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Random Usenrame
	//||------------------------------------------------------------------------------------------------||

	randomUsername, err := actions.GenerateUsername()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("USERNAME_GEN_FAILED"))
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Record
	//||------------------------------------------------------------------------------------------------||

	fmt.Println("Creating account for:", session.Identifier)
	fmt.Println(actions.GenerateIdentifierHash(session.Identifier))
	account := db.ModelAccount{}
	account.ID = dbAccount.ID
	account.Identifier = actions.GenerateIdentifierHash(session.Identifier)
	account.Username = randomUsername
	account.Salt = saltHash
	account.Level = 1
	account.Status = app.Constants("AccountStatus").Code("Active")
	db.AuthDB().Save(&account)

	//||------------------------------------------------------------------------------------------------||
	//|| Refetch the User Data
	//||------------------------------------------------------------------------------------------------||

	updatedAccount, err := db.GetAccountByID(fmt.Sprintf("%d", account.ID))
	if err != nil || updatedAccount == nil {
		responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("ACCOUNT_LOOKUP_FAILED"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Session
	//||------------------------------------------------------------------------------------------------||

	sessionToken, err := actions.SessionCreate(updatedAccount.Identifier, updatedAccount)
	if err != nil || sessionToken == "" {
		responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("SESSION_GEN_FAILED"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Write the Session Cookie
	//||------------------------------------------------------------------------------------------------||

	actions.WriteSessionCookie(w, sessionToken)

	//||------------------------------------------------------------------------------------------------||
	//|| Delete the Old Session Cookie
	//||------------------------------------------------------------------------------------------------||

	if cookie.Value != "" && cookie.Value != sessionToken {
		_ = actions.DeleteSession(cookie.Value)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Success
	//||------------------------------------------------------------------------------------------------||

	responses.Success(w, http.StatusOK, authCompleteResponse{
		Message: "",
		Next:    "/members",
	})
}
