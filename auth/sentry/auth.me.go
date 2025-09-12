package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"base/verify"
	"encoding/json"
	"net/http"

	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/auth/actions"
)

//||------------------------------------------------------------------------------------------------||
//|| Handler
//||------------------------------------------------------------------------------------------------||

type authMeResponse struct {
	ID       int64           `json:"id"`
	Status   string          `json:"status"`
	Type     string          `json:"type"`
	Email    string          `json:"email"`
	Username string          `json:"username"`
	Level    int             `json:"level"`
	Security int             `json:"security"`
	Created  int64           `json:"created"`
	Expires  int64           `json:"expires"`
	Identity verify.Identity `json:"identity"`
}

//||------------------------------------------------------------------------------------------------||
//|| Handler
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
	//|| Identity
	//||------------------------------------------------------------------------------------------------||

	var ident verify.Identity
	if len(account.Identity) > 0 {
		if err := json.Unmarshal([]byte(account.Identity), &ident); err != nil {
			responses.Error(w, http.StatusInternalServerError, "Failed to parse identity")
			return
		}
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Return Responses
	//||------------------------------------------------------------------------------------------------||

	response := authMeResponse{
		ID:       account.ID,
		Status:   account.Status,
		Type:     account.Type,
		Email:    session.Identifier,
		Username: account.Username,
		Level:    account.Level,
		Security: account.Security,
		Created:  session.Created,
		Expires:  session.Expires,
		Identity: ident,
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Return Response
	//||------------------------------------------------------------------------------------------------||

	responses.Success(w, http.StatusOK, response)

}
