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

	valid := validate.ValidatePhoneOrEmail(identifier)
	if valid != nil {
		responses.Error(w, http.StatusBadRequest, valid.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Validate Password
	//||------------------------------------------------------------------------------------------------||

	vp := validate.IsValidPassword(password)
	if vp != nil {
		responses.Error(w, http.StatusBadRequest, vp.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Call the Acccount Return
	//||------------------------------------------------------------------------------------------------||

	hashedIdentifier := actions.GenerateIdentifierHash(identifier)
	account, err := db.GetAccountByIdentifier(hashedIdentifier)
	if err != nil || account == nil {
		actions.VerifyPassword("NOT_ACTUALLY_VERIFYING", "THIS IS JUST TO_PREVENT_TIMING_ATTACKS", "AND_DATA_LEAKAGE")
		responses.Error(w, http.StatusUnauthorized, app.Err("Auth").Code("INVALID_CREDENTIALS"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Check the Password
	//||------------------------------------------------------------------------------------------------||

	if !actions.VerifyPassword(account.Salt, password, account.Password) {
		responses.Error(w, http.StatusUnauthorized, app.Err("Auth").Code("INVALID_CREDENTIALS"))
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

	responses.Error(w, http.StatusUnauthorized, app.Err("Auth").Code("SESSION_GEN_FAILED"))

}
