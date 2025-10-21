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
	"github.com/ralphferrara/aria/log"
	"github.com/ralphferrara/aria/responses"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/auth/db"
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

	app.Log.Data("[TwoFactor] Incoming -> token:", token, " code:", code)

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

	app.Log.Data("[TwoFactor] Record:")
	log.PrettyPrint(record)

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
	fmt.Println("Successfully validated Record:", record.Type)

	//||------------------------------------------------------------------------------------------------||
	//|| Get the Hashed Email
	//||------------------------------------------------------------------------------------------------||

	fmt.Println("Identifier:", record.Identifier)
	hashedIdentifier := actions.GenerateIdentifierHash(record.Identifier)
	fmt.Println("Hashed Identifier:", hashedIdentifier)

	//||------------------------------------------------------------------------------------------------||
	//|| PASSWORD RESET VERIFICATION
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

	account, _ := db.GetAccountByIdentifier(hashedIdentifier)
	fmt.Println("Account Lookup:", hashedIdentifier)
	fmt.Println(account)

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
		account, uErr = db.CreateAccount(&created)
		if uErr != nil || account == nil {
			responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("ACCOUNT_CREATE_FAILED"))
			return
		}
		//||------------------------------------------------------------------------------------------------||
		//|| Update Status to Verified
		//||------------------------------------------------------------------------------------------------||
		err = db.UpdateStatusVerified(uint(account.ID))
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, app.Err("Auth").Code("ACCOUNT_CREATE_FAILED"))
			return
		}
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Record
	//||------------------------------------------------------------------------------------------------||

	newToken, err := actions.SessionCreate(record.Identifier, account)
	fmt.Println("Session Create:", newToken, err)
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
