package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"base/db/models"
	"base/verify"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/auth/types"
)

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

type responseTwoFactor struct {
	Next    string `json:"next"`
	Message string `json:"message"`
}

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

func TwoFactorHandler(w http.ResponseWriter, r *http.Request) {
	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||
	token := r.FormValue("token")
	code := r.FormValue("code")
	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||
	fmt.Println("[TwoFactor] Incoming -> token:", token, " code:", code)
	//||------------------------------------------------------------------------------------------------||
	//|| Basic Validation
	//||------------------------------------------------------------------------------------------------||
	if token == "" || code == "" {
		responses.Error(w, http.StatusBadRequest, "Missing or invalid code/token/type")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Get Record from Redis
	//||------------------------------------------------------------------------------------------------||
	val, err := app.CacheRedis["auth"].Get(actions.TwoFactorCacheCode(token))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid or expired token")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Convert to struct
	//||------------------------------------------------------------------------------------------------||
	var record types.TwoFactorVerification
	if err := json.Unmarshal([]byte(val), &record); err != nil {
		responses.Error(w, http.StatusInternalServerError, "Invalid stored record")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Get Record from Redis
	//||------------------------------------------------------------------------------------------------||
	if record.Attempts >= 5 {
		app.CacheRedis["auth"].Del(fmt.Sprintf("verify:%s", token))
		responses.Error(w, http.StatusTooManyRequests, "Too many attempts")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Get Record from Redis
	//||------------------------------------------------------------------------------------------------||
	if code != record.Code {
		record.Attempts++
		newData, _ := json.Marshal(record)
		app.CacheRedis["auth"].Set(fmt.Sprintf("verify:%s", token), newData, time.Until(record.Expires))
		responses.Error(w, http.StatusUnauthorized, "Invalid code")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Check if the token is expired
	//||------------------------------------------------------------------------------------------------||
	if time.Now().After(record.Expires) {
		app.CacheRedis["auth"].Del(fmt.Sprintf("verify:%s", token))
		responses.Error(w, http.StatusBadRequest, "Token expired")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Success! Delete
	//||------------------------------------------------------------------------------------------------||
	app.CacheRedis["auth"].Del(fmt.Sprintf("verify:%s", token))
	//||------------------------------------------------------------------------------------------------||
	//|| Get the Hashed Email
	//||------------------------------------------------------------------------------------------------||
	hashedEmail := actions.GenerateEmailHash(record.Identifier)
	//||------------------------------------------------------------------------------------------------||
	//|| Get the Account Record
	//||------------------------------------------------------------------------------------------------||
	var account models.Account
	if err := app.SQLDB["main"].DB.Where("account_email = ?", hashedEmail).First(&account).Error; err != nil {
		fmt.Println("[Session] Account not found for email:", record.Identifier, " creating new account")
	} else {
		existsToken, err := actions.SessionCreate(record.Identifier, &account)
		if err == nil {
			actions.WriteSessionCookie(w, existsToken)
			responses.Success(w, http.StatusOK, responseTwoFactor{
				Message: "Two-factor authentication successful. Redirecting to /members",
				Next:    "/members/",
			})
			return
		}
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Identity
	//||------------------------------------------------------------------------------------------------||

	identity := verify.Identity{}
	identityJSON, err := json.Marshal(identity)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to marshal identity: "+err.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Struct
	//||------------------------------------------------------------------------------------------------||
	account.Type = record.Type
	account.Email = hashedEmail
	account.Status = "VERF"
	account.Level = 1
	account.Identity = string(identityJSON)
	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Record
	//||------------------------------------------------------------------------------------------------||
	if err := app.SQLDB["main"].DB.Create(&account).Error; err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to create account")
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Record
	//||------------------------------------------------------------------------------------------------||
	newToken, err := actions.SessionCreate(record.Identifier, &account)
	if err == nil {
		actions.WriteSessionCookie(w, newToken)
		responses.Success(w, http.StatusOK, responseTwoFactor{
			Message: "Account Created",
			Next:    "/complete",
		})
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Record
	//||------------------------------------------------------------------------------------------------||

	responses.Error(w, http.StatusInternalServerError, "Failed to create account - Unknown")
}
