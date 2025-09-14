package db

import (
	"fmt"

	"github.com/ralphferrara/aria/app"
)

func UpdatePassword(accountID int64, passwordHash string) error {
	err := AuthDB().Model(&ModelAccount{}).
		Where("id_account = ?", accountID).
		Updates(map[string]any{
			"account_password": passwordHash,
		}).Error

	if err != nil {
		fmt.Println("Error updating password:", err)
		return app.Err("Auth").Error("PASSWORD_UPDATE_FAILED")
	}
	return nil
}
