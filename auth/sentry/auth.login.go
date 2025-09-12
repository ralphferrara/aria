package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"net/http"

	"base/db/abstract"

	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/base/validate"
	"github.com/ralphferrara/aria/responses"
)

//||------------------------------------------------------------------------------------------------||
//|| Response
//||------------------------------------------------------------------------------------------------||

type loginResponse struct {
	Next string `json:"next"`
}

//||------------------------------------------------------------------------------------------------||
//|| Handler :: Processes the Login Request
//||------------------------------------------------------------------------------------------------||

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")
	captcha := r.FormValue("captcha")

	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||

	fmt.Println("[Login]")
	fmt.Println("identifier:", identifier)
	fmt.Println("password:", password)
	fmt.Println("captcha:", captcha)

	//||------------------------------------------------------------------------------------------------||
	//|| Validate
	//||------------------------------------------------------------------------------------------------||

	if !validate.IsValidEmail(identifier) {
		responses.Error(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Validate Password
	//||------------------------------------------------------------------------------------------------||

	if password == "" {
		responses.Error(w, http.StatusBadRequest, "Password is required")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Call the Acccount Return
	//||------------------------------------------------------------------------------------------------||

	hashedEmail := actions.GenerateEmailHash(identifier)
	fmt.Println(identifier)
	fmt.Println("Hashed Account:", hashedEmail)
	account, err := abstract.GetAccountByEmail(hashedEmail)
	if err != nil || account == nil {
		responses.Error(w, http.StatusUnauthorized, "Invalid email or password1")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Check the Password
	//||------------------------------------------------------------------------------------------------||

	if !actions.VerifyPassword(account.Salt, password, account.Password) {
		responses.Error(w, http.StatusUnauthorized, "Invalid email or password2")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Get the Account Record and redirect if completed
	//||------------------------------------------------------------------------------------------------||

	loginToken, err := actions.SessionCreate(identifier, account)
	if err == nil {
		actions.WriteSessionCookie(w, loginToken)
		responses.Success(w, http.StatusOK, loginResponse{
			Next: "/members/",
		})
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Store
	//||------------------------------------------------------------------------------------------------||

	responses.Error(w, http.StatusUnauthorized, "Error Writing Session Cookie, please try again later")

}
