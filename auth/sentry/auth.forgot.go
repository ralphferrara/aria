package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/types"
	"github.com/ralphferrara/aria/base/random"
	"github.com/ralphferrara/aria/base/validate"
)

//||------------------------------------------------------------------------------------------------||
//|| Response
//||------------------------------------------------------------------------------------------------||

type responseForgotPassword struct {
	Token      string `json:"token"`
	Identifier string `json:"identifier"`
}

//||------------------------------------------------------------------------------------------------||
//|| Handler :: Initiates Forgot Password Flow
//||------------------------------------------------------------------------------------------------||

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||

	identifier := r.FormValue("identifier")
	idType := validate.IsEmailOrPhone(identifier)

	//||------------------------------------------------------------------------------------------------||
	//|| Validate Email
	//||------------------------------------------------------------------------------------------------||

	if idType == "email" && !validate.IsValidEmail(identifier) {
		responses.Error(w, http.StatusBadRequest, app.Err("Auth").Code("INVALID_EMAIL"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Validate Phone
	//||------------------------------------------------------------------------------------------------||

	if idType == "phone" && !validate.IsValidPhone(identifier) {
		responses.Error(w, http.StatusBadRequest, app.Err("Auth").Code("INVALID_PHONE"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Generate Reset Token & Code
	//||------------------------------------------------------------------------------------------------||

	keyEncoded := random.UUIDString()
	keyCode := random.NumberString(6)

	//||------------------------------------------------------------------------------------------------||
	//|| Create Forgot Password Record
	//||------------------------------------------------------------------------------------------------||

	record := types.TwoFactorVerification{
		Code:       keyCode,
		Key:        keyEncoded,
		Identifier: identifier,
		Type:       app.Constants("TwoFactorType").Code("Reset"),
		Attempts:   0,
		Created:    time.Now(),
		Expires:    time.Now().Add(15 * time.Minute),
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Serialize Record for Storage
	//||------------------------------------------------------------------------------------------------||

	data, err := json.Marshal(record)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to serialize reset request")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Save Reset Request in Redis (15 min expiry)
	//||------------------------------------------------------------------------------------------------||

	err = app.CacheRedis["auth"].Set("reset::"+keyEncoded, data, 15*time.Minute)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to store reset request")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Return Success
	//||------------------------------------------------------------------------------------------------||

	responses.Success(w, http.StatusOK, responseForgotPassword{
		Token:      keyEncoded,
		Identifier: identifier,
	})
}
