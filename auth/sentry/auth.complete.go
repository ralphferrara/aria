package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"api/send"
	"base/db/abstract"
	"base/db/models"
	"base/verify"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/actions"
	"github.com/ralphferrara/aria/base/crypto"
	"github.com/ralphferrara/aria/base/encrypt"
	"github.com/ralphferrara/aria/base/validate"
	"github.com/ralphferrara/aria/responses"
)

//||------------------------------------------------------------------------------------------------||
//|| Response
//||------------------------------------------------------------------------------------------------||

type authCompleteResponse struct {
	Message string `json:"message"`
	Next    string `json:"next"`
}

//||------------------------------------------------------------------------------------------------||
//|| Handler
//||------------------------------------------------------------------------------------------------||

func CompleteHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| DB Account
	//||------------------------------------------------------------------------------------------------||

	cookie, dbAccount, session, err := actions.LoadSessionAccount(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Check Account Status
	//||------------------------------------------------------------------------------------------------||

	if dbAccount.Status != verify.StatusVerified.String() {
		responses.Error(w, http.StatusForbidden, "Account is already created")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||

	password := r.FormValue("password")
	rawEncrypt := r.FormValue("encryptionLevel")
	privateKeyInput := r.FormValue("privateKey")
	publicKeyInput := r.FormValue("publicKey")
	wordListJSON := r.FormValue("wordList")

	//||------------------------------------------------------------------------------------------------||
	//||
	//|| Sanitize and Validate
	//|| Also generate the private/public key if needed
	//||
	//||------------------------------------------------------------------------------------------------||

	if password == "" || len(password) < 8 {
		responses.Error(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Validate Encryption Level
	//||------------------------------------------------------------------------------------------------||

	encryptionLevel, err := strconv.Atoi(rawEncrypt)
	if err != nil || encryptionLevel < 1 || encryptionLevel > 3 {
		responses.Error(w, http.StatusBadRequest, "Invalid or missing encryption level")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Validate and Generate Private/Public Key
	//||------------------------------------------------------------------------------------------------||

	var privateKey, publicKey string

	//||------------------------------------------------------------------------------------------------||
	//|| Level 1 - We handle the keys
	//||------------------------------------------------------------------------------------------------||

	if encryptionLevel == 1 {
		genPrivateKey, genPublicKey, err := crypto.GenerateKeyPair()
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, "Failed to generate keys")
			return
		}
		privateKey = genPrivateKey
		publicKey = genPublicKey
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Level 2 - BIPList
	//||------------------------------------------------------------------------------------------------||

	var BIPList []string

	if encryptionLevel == 2 {
		BIPList, err := validate.ValidateBIP39(wordListJSON)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "Invalid BIP39 word list: "+err.Error())
			return
		}
		genPrivate, genPublic, err := crypto.GenerateBIP39Keys(BIPList)
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, "Failed to generate BIP39 keys")
			return
		}
		privateKey = genPrivate
		publicKey = genPublic
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Level 3 requires both keys
	//||------------------------------------------------------------------------------------------------||

	if encryptionLevel == 3 {
		err := validate.ValidateKeyPair(privateKeyInput, publicKeyInput)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "Invalid key pair: "+err.Error())
			return
		}
		privateKey = privateKeyInput
		publicKey = publicKeyInput
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Generate the Private Key Hash
	//||------------------------------------------------------------------------------------------------||

	privateKeyHash, err := encrypt.GenerateCheckKey(privateKey)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to generate private key hash")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Password/Salt
	//||------------------------------------------------------------------------------------------------||

	passwordHash, saltHash := actions.GeneratePassword(password)
	if passwordHash == "" {
		responses.Error(w, http.StatusBadRequest, "Could not generate password")
	}

	//||------------------------------------------------------------------------------------------------||
	//||
	//|| Contact the User with the keys if needed
	//||
	//||------------------------------------------------------------------------------------------------||

	if encryptionLevel == 1 {
		_ = send.EmailPrivateKeyToUser(session.Identifier, privateKey)
	}

	if encryptionLevel == 2 {
		_ = send.EmailBIPListToUser(session.Identifier, BIPList, privateKey)
	}

	if encryptionLevel == 3 {
		_ = send.EmailPrivateKeyToUser(session.Identifier, privateKey)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Random Usenrame
	//||------------------------------------------------------------------------------------------------||

	randomUsername, err := actions.GenerateUsername()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to generate username")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Email Verification
	//||------------------------------------------------------------------------------------------------||

	verifyRecord, err := verify.Create(verify.DataTypeMAIL, dbAccount.ID, app.Storages["verifications"], app.SQLDB["main"], privateKey, publicKey)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to initialize verification: "+err.Error())
		return
	}
	verifyRecord.UpdateStatusVerified("TWOFACTOR") // Automatic Moderator

	//||------------------------------------------------------------------------------------------------||
	//|| Identity
	//||------------------------------------------------------------------------------------------------||

	verifyRecord.UpdateVerification(verify.DataTypeMAIL, verifyRecord.Display, verifyRecord.UUID)
	identity := verify.Identity{}
	identityJSON, err := json.Marshal(identity)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to marshal identity: "+err.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Account Record
	//||------------------------------------------------------------------------------------------------||

	fmt.Println("Creating account for:", session.Identifier)
	fmt.Println(actions.GenerateEmailHash(session.Identifier))
	account := models.Account{}
	account.ID = dbAccount.ID
	account.Username = dbAccount.Username
	account.Type = dbAccount.Type
	account.Email = actions.GenerateEmailHash(session.Identifier)
	account.Username = randomUsername
	account.Public = publicKey
	account.Password = passwordHash
	account.PrivateHash = privateKeyHash
	account.Salt = saltHash
	account.Level = 1
	account.Status = "ACTV"
	account.Security = encryptionLevel // Default security level
	account.Identity = string(identityJSON)
	account.Private = privateKey
	app.SQLDB["main"].DB.Save(&account)

	//||------------------------------------------------------------------------------------------------||
	//|| Refetch the User Data
	//||------------------------------------------------------------------------------------------------||

	updatedAccount, err := abstract.GetAccountByID(fmt.Sprintf("%d", account.ID))
	if err != nil || updatedAccount == nil {
		responses.Error(w, http.StatusInternalServerError, "Could not re-fetch account after update")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create the Session
	//||------------------------------------------------------------------------------------------------||

	sessionToken, err := actions.SessionCreate(updatedAccount.Email, updatedAccount)
	if err != nil || sessionToken == "" {
		responses.Error(w, http.StatusInternalServerError, "Failed to create session: "+err.Error())
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Write the Session Cookie
	//||------------------------------------------------------------------------------------------------||

	actions.WriteSessionCookie(w, sessionToken)

	//||------------------------------------------------------------------------------------------------||
	//|| Delete the Old Session Cookie
	//||------------------------------------------------------------------------------------------------||

	if cookie.Value != "" && cookie.Value != sessionToken {
		_ = actions.DeleteSession(cookie.Value)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Success
	//||------------------------------------------------------------------------------------------------||

	responses.Success(w, http.StatusOK, authCompleteResponse{
		Message: "Signup complete",
		Next:    "/members",
	})
}
