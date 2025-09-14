package db

import "gorm.io/gorm"

//||------------------------------------------------------------------------------------------------||
//|| CreateAccount
//||------------------------------------------------------------------------------------------------||

func CreateAccount(account *ModelAccount) (*ModelAccount, error) {
	result := AuthDB().Create(account)
	if result.Error != nil {
		return nil, result.Error
	}

	err := AuthDB().First(account, account.ID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return account, nil
}
