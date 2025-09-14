package db

import (
	"gorm.io/gorm"
)

//||------------------------------------------------------------------------------------------------||
//|| Get Account Based on ID
//||------------------------------------------------------------------------------------------------||

func GetAccountByID(id string) (*ModelAccount, error) {
	var account ModelAccount

	result := AuthDB().First(&account, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &account, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Get Account Based on Email
//||------------------------------------------------------------------------------------------------||

func GetAccountByEmail(hashedIdentifier string) (*ModelAccount, error) {
	var account ModelAccount

	result := AuthDB().Where("account_identifier = ?", hashedIdentifier).Limit(1).Find(&account)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &account, nil
}
