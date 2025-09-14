package db

import (
	"fmt"

	"github.com/ralphferrara/aria/app"
)

func GetAccountByIdentifier(identifier string) (*ModelAccount, error) {
	account := &ModelAccount{}
	err := AuthDB().Where("account_identifier = ?", identifier).First(account).Error
	if err != nil {
		fmt.Println("[DB] GetAccountByIdentifier error:", err)
		return nil, app.Err("Auth").Error("ACCOUNT_NOT_FOUND")
	}
	return account, nil
}
