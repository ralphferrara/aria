package actions

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"base/db/abstract"
	"base/db/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/ralphferrara/aria/auth/types"
)

//||------------------------------------------------------------------------------------------------||
//|| Load Session Account
//||------------------------------------------------------------------------------------------------||

func LoadSessionAccount(r *http.Request) (http.Cookie, models.Account, types.SessionRecord, error) {

	//||------------------------------------------------------------------------------------------------||
	//|| Get the Session Cookie
	//||------------------------------------------------------------------------------------------------||

	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		return http.Cookie{}, models.Account{}, types.SessionRecord{}, errors.New("missing session cookie")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Get Session
	//||------------------------------------------------------------------------------------------------||

	session, err := FetchSession(cookie.Value)
	if err != nil {
		return *cookie, models.Account{}, types.SessionRecord{}, errors.New("could not retrieve session")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Get Database Account
	//||------------------------------------------------------------------------------------------------||

	account, err := abstract.GetAccountByID(fmt.Sprintf("%d", session.ID))
	if err != nil || account == nil {
		return *cookie, models.Account{}, types.SessionRecord{}, errors.New("could not retrieve account")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Done
	//||------------------------------------------------------------------------------------------------||

	return *cookie, *account, session, nil
}
