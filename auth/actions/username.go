//||------------------------------------------------------------------------------------------------||
//|| GenerateUsername: Creates a Unique Random Username
//||------------------------------------------------------------------------------------------------||

package actions

import (
	"fmt"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/db"
	"github.com/ralphferrara/aria/base/bip39"
	"github.com/ralphferrara/aria/base/random"
)

//||------------------------------------------------------------------------------------------------||
//|| GenerateUsername
//||------------------------------------------------------------------------------------------------||

func GenerateUsername() (string, error) {
	for range make([]struct{}, 5) {
		word1 := bip39.RandomBIP39Word()
		word2 := bip39.RandomBIP39Word()
		word3 := bip39.RandomBIP39Word()
		word4 := random.RandomString(3)
		//||------------------------------------------------------------------------------------------------||
		//|| Generate Random Username
		//||------------------------------------------------------------------------------------------------||
		username := fmt.Sprintf("%s_%s_%s_%s", word1, word2, word3, word4)
		exists := db.CheckUsernameExists(username)
		if !exists {
			return username, nil
		}
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Success
	//||------------------------------------------------------------------------------------------------||
	return "", app.Err("Auth").Error("USERNAME_GEN_FAILED")
}
