package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ralphferrara/aria/base/random"
	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/auth/db"
	"github.com/ralphferrara/aria/auth/setup"
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
		responses.Error(w, http.StatusBadRequest, app.Err("Auth").Code("TF_MISSING_CODE"))
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Get Record from Redis
	//||------------------------------------------------------------------------------------------------||
	val, err := app.CacheRedis["auth"].Get(actions.TwoFactorCacheCode(token))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, app.Err("Auth").Code("TF_INVALID_TOKEN"))
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Convert to struct
	//||------------------------------------------------------------------------------------------------||
	var record types.TwoFactorVerification
	if err := json.Unmarshal([]byte(val), &record); err != nil {
		responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("TF_INVALID_RECORD"))
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Get Record from Redis
	//||------------------------------------------------------------------------------------------------||
	if record.Attempts >= 5 {
		app.CacheRedis["auth"].Del(fmt.Sprintf("verify:%s", token))
		responses.Error(w, http.StatusTooManyRequests, app.Err("Auth").Code("TF_TOO_MANY_ATTEMPTS"))
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Get Record from Redis
	//||------------------------------------------------------------------------------------------------||
	if code != record.Code {
		record.Attempts++
		newData, _ := json.Marshal(record)
		app.CacheRedis["auth"].Set(fmt.Sprintf("verify:%s", token), newData, time.Until(record.Expires))
		responses.Error(w, http.StatusUnauthorized, app.Err("Auth").Code("TF_CODE_MISMATCH"))
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Check if the token is expired
	//||------------------------------------------------------------------------------------------------||
	if time.Now().After(record.Expires) {
		app.CacheRedis["auth"].Del(fmt.Sprintf("verify:%s", token))
		responses.Error(w, http.StatusBadRequest, app.Err("Auth").Code("TF_TOKEN_EXPIRED"))
		return
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Success! Delete
	//||------------------------------------------------------------------------------------------------||
	app.CacheRedis["auth"].Del(fmt.Sprintf("verify:%s", token))
	//||------------------------------------------------------------------------------------------------||
	//|| Get the Hashed Email
	//||------------------------------------------------------------------------------------------------||
	hashedIdentifier := actions.GenerateIdentifierHash(record.Identifier)

	//||------------------------------------------------------------------------------------------------||
	//|| Handle Password Reset Verification
	//||------------------------------------------------------------------------------------------------||

	if record.Type == app.Constants("TwoFactorType").Code("Reset") {

		//||------------------------------------------------------------------------------------------------||
		//|| Lookup the Account
		//||------------------------------------------------------------------------------------------------||

		account, err := db.GetAccountByIdentifier(hashedIdentifier)
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("ACCOUNT_NOT_FOUND"))
			return
		}

		//||------------------------------------------------------------------------------------------------||
		//|| Create the Session and redirect to reset password
		//||------------------------------------------------------------------------------------------------||

		existsToken, err := actions.SessionCreate(account.Identifier, account)
		if err == nil {
			actions.WriteSessionCookie(w, existsToken)
			responses.Success(w, http.StatusOK, responseTwoFactor{
				Message: "OK",
				Next:    "/reset",
			})
			return
		}

		responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("SESSION_GEN_FAILED"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Account Creation - NOT RESET
	//||------------------------------------------------------------------------------------------------||

	account, aErr := db.GetAccountByIdentifier(hashedIdentifier)
	if aErr != nil {
		responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("ACCOUNT_LOOKUP_FAILED"))
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account
	//||------------------------------------------------------------------------------------------------||

	nextPage := "/complete"
	if account == nil {
		nextPage = "/members/"
		//||------------------------------------------------------------------------------------------------||
		//|| Create Username
		//||------------------------------------------------------------------------------------------------||
		username, uErr := actions.GenerateUsername()
		if uErr != nil {
			responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("USERNAME_GEN_FAILED"))
			return
		}
		//||------------------------------------------------------------------------------------------------||
		//|| Create Account Record
		//||------------------------------------------------------------------------------------------------||
		created := db.ModelAccount{}
		created.Identifier = hashedIdentifier
		created.Username = username
		created.Salt = random.RandomString(32)
		created.Status = app.Constants("AccountStatus").Code("Pending")
		created.Level = 1
		account, aErr = db.CreateAccount(&created)
		if aErr != nil || account == nil {
			responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("ACCOUNT_CREATE_FAILED"))
			return
		}
		create := setup.Setup.Functions.OnAccountCreation(r, account.ID)
		if create != nil {
			responses.Error(w, http.StatusInternalServerError, create.Error())
			return
		}
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Record
	//||------------------------------------------------------------------------------------------------||

	newToken, err := actions.SessionCreate(record.Identifier, account)
	if err == nil {
		actions.WriteSessionCookie(w, newToken)
		responses.Success(w, http.StatusOK, responseTwoFactor{
			Message: "OK",
			Next:    nextPage,
		})
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Record
	//||------------------------------------------------------------------------------------------------||

	responses.Error(w, http.StatusInternalServerError, "Failed to create account - Unknown")
}
