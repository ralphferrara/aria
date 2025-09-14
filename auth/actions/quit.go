package actions

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ralphferrara/aria/auth/setup"
)

//||------------------------------------------------------------------------------------------------||
//|| GenerateQuitToken – creates an HMAC hash from AccountPrivateHash
//||------------------------------------------------------------------------------------------------||

func GenerateQuitToken(accountPassword string) string {
	secret := []byte(setup.Setup.CSRF)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(accountPassword))
	return hex.EncodeToString(h.Sum(nil))
}
