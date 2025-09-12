package sentry

import (
	"fmt"
	"net/http"

	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/auth/actions"
)

//||------------------------------------------------------------------------------------------------||
//|| Response
//||------------------------------------------------------------------------------------------------||

type logoutResponse struct {
	Message string `json:"message"`
	Next    string `json:"next"`
}

//||------------------------------------------------------------------------------------------------||
//|| Logout Handler
//||------------------------------------------------------------------------------------------------||

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| Get the Session Cookie
	//||------------------------------------------------------------------------------------------------||

	cookie, _, _, err := actions.LoadSessionAccount(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Delete Session from Redis
	//||------------------------------------------------------------------------------------------------||

	err = actions.DeleteSession(cookie.Value)
	if err != nil {
		fmt.Printf("[Logout] Failed to delete session %s from Redis: %v\n", cookie.Value, err)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Clear the Session Cookie from the browser
	//||------------------------------------------------------------------------------------------------||

	actions.ClearSessionCookie(w)

	//||------------------------------------------------------------------------------------------------||
	//|| Success Response
	//||------------------------------------------------------------------------------------------------||

	responses.Success(w, http.StatusOK, logoutResponse{
		Message: "Logged out successfully",
		Next:    "/login",
	})
}
