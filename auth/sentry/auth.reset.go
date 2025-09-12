package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"base/db/abstract"
	"base/db/models"
	"fmt"
	"net/http"

	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/actions"
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

	account, err := abstract.GetAccountByID(fmt.Sprintf("%d", session.ID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to fetch account")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Parse Password Fields
	//||------------------------------------------------------------------------------------------------||

	password := r.FormValue("password")
	confirm := r.FormValue("confirmPassword")

	fmt.Println("Resetting password for account ID:", session.ID)
	fmt.Println("Password:", password)
	fmt.Println("Confirm Password:", confirm)

	if password == "" || confirm == "" {
		responses.Error(w, http.StatusBadRequest, "Password fields cannot be empty")
		return
	}
	if len(password) < 8 {
		responses.Error(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}
	if password != confirm {
		responses.Error(w, http.StatusBadRequest, "Passwords do not match")
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

	err = app.SQLDB["main"].DB.Model(&models.Account{}).
		Where("id_account = ?", session.ID).
		Updates(map[string]any{
			"account_password": passwordHash,
		}).Error

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Success
	//||------------------------------------------------------------------------------------------------||
	responses.Success(w, http.StatusOK, map[string]any{
		"message": "Password successfully updated",
	})
}
