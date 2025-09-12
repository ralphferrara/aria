package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/auth/types"
	"github.com/ralphferrara/aria/base/random"
	"github.com/ralphferrara/aria/base/validate"
)

//||------------------------------------------------------------------------------------------------||
//|| Response
//||------------------------------------------------------------------------------------------------||

type signupResponse struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

//||------------------------------------------------------------------------------------------------||
//|| Handler :: Processes the Signup Request
//||------------------------------------------------------------------------------------------------||

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||
	email := r.FormValue("email")
	accountType := r.FormValue("type")
	//||------------------------------------------------------------------------------------------------||
	//|| Validate
	//||------------------------------------------------------------------------------------------------||
	if !validate.IsValidEmail(email) {
		responses.Error(w, http.StatusBadRequest, "Invalid email")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Generate the Values
	//||------------------------------------------------------------------------------------------------||
	key := random.UUIDString()
	code := random.NumberString(6)
	//||------------------------------------------------------------------------------------------------||
	//|| Create record
	//||------------------------------------------------------------------------------------------------||
	record := types.TwoFactorVerification{
		Code:       code,
		Key:        key,
		Identifier: email,
		Type:       accountType,
		Attempts:   0,
		Created:    time.Now(),
		Expires:    time.Now().Add(15 * time.Minute),
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Serialize to Save
	//||------------------------------------------------------------------------------------------------||
	data, err := json.Marshal(record)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to serialize")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Save to Redis with expiry
	//||------------------------------------------------------------------------------------------------||

	err = app.CacheRedis["auth"].Set(actions.TwoFactorCacheCode(key), data, 15*time.Minute)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to cache verification")
		return
	}
	fmt.Printf("✅ TwoFactor %s :: key=%s code=%s\n", accountType, key, code)

	//||------------------------------------------------------------------------------------------------||
	//|| Store
	//||------------------------------------------------------------------------------------------------||

	responses.Success(w, http.StatusOK, signupResponse{
		Token: key,
		Type:  accountType,
	})
}
