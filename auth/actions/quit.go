package actions

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ralphferrara/aria/auth"
)

//||------------------------------------------------------------------------------------------------||
//|| GenerateQuitToken – creates an HMAC hash from AccountPrivateHash
//||------------------------------------------------------------------------------------------------||

func GenerateQuitToken(accountPrivateHash string) string {
	secret := []byte(auth.Setup.CSRF)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(accountPrivateHash))
	return hex.EncodeToString(h.Sum(nil))
}
