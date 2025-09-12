package defauth

import (
	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth"
	"github.com/ralphferrara/aria/auth/types"
	"github.com/ralphferrara/aria/base/encrypt"
)

//||------------------------------------------------------------------------------------------------||
//|| Get the Account Record
//||------------------------------------------------------------------------------------------------||

func FetchAccountByIdentifier(identifier string) types.AccountAuth {

	//||------------------------------------------------------------------------------------------------||
	//|| It's Set
	//||------------------------------------------------------------------------------------------------||

	if auth.Setup.FetchAccountByIdentifier != nil {
		return auth.Setup.FetchAccountByIdentifier(identifier)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Default - Hashed Email
	//||------------------------------------------------------------------------------------------------||

	hashedIdentifier := encrypt.KeyHash(Setup.Pepper, identifier)

	//||------------------------------------------------------------------------------------------------||
	//|| Get the Account Record
	//||------------------------------------------------------------------------------------------------||

	var account types.AccountAuth
	record, err := app.SQLDB["main"].DB.Where("account_email = ?", hashedIdentifier).First(&account).Error
	if err != nil {
		return types.AccountAuth{
			Identifier : "INVALID-IDENTIFIER"
			HashedPassword : random.GenerateString(64)
			Pepper : Setup.Pepper
			Salt : "DUMMY-PREVENT-TIMING-DUMMY-PREVENT-TIMING-DUMMY-PREVENT-TIMING"
			Level : 0
		}
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Record
	//||------------------------------------------------------------------------------------------------||

	return types.AccountAuth{ 
		Identifier: record.Identifier,
		HashedPassword: record.Password,
		Salt: record.Salt,
		Pepper: Setup.Pepper,
		Level: record.Level,
	}
}
