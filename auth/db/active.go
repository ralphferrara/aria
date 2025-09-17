package db

import (
	"fmt"

	"github.com/ralphferrara/aria/app"
)

func UpdateAccountActive(accountID int64) error {
	err := AuthDB().Model(&ModelAccount{}).
		Where("id_account = ?", accountID).
		Updates(map[string]any{
			"account_status": app.Constants("AccountStatus").Code("Active"),
		}).Error

	if err != nil {
		fmt.Println("Error updating account status to active :", err)
		return app.Err("Auth").Error("PASSWORD_UPDATE_FAILED")
	}
	return nil
}
