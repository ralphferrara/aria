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
	"github.com/ralphferrara/aria/auth/types"
	"github.com/ralphferrara/aria/base/random"
	"github.com/ralphferrara/aria/base/validate"
)

//||------------------------------------------------------------------------------------------------||
//|| Response
//||------------------------------------------------------------------------------------------------||

type responseForgotPassword struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

//||------------------------------------------------------------------------------------------------||
//|| Handler :: Initiates Forgot Password Flow
//||------------------------------------------------------------------------------------------------||

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||

	email := r.FormValue("email")

	//||------------------------------------------------------------------------------------------------||
	//|| Validate Email
	//||------------------------------------------------------------------------------------------------||

	if !validate.IsValidEmail(email) {
		responses.Error(w, http.StatusBadRequest, "Invalid email address")
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
		Identifier: email,
		Type:       "RESET",
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

	fmt.Printf("✅ ForgotPassword :: email=%s key=%s code=%s\n", email, keyEncoded, keyCode)

	//||------------------------------------------------------------------------------------------------||
	//|| Return Success
	//||------------------------------------------------------------------------------------------------||

	responses.Success(w, http.StatusOK, responseForgotPassword{
		Token: keyEncoded,
		Email: email,
	})
}
