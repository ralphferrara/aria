package db

import "github.com/ralphferrara/aria/auth/setup"

//||------------------------------------------------------------------------------------------------||
//|| DB ModelEmail Verification
//||------------------------------------------------------------------------------------------------||

type ModelAccount struct {
	ID         int64  `gorm:"column:id_account;primaryKey;autoIncrement"`
	Salt       string `gorm:"column:account_salt;size:256"`
	Username   string `gorm:"column:account_username;size:64;index:idx_accounts_account_username"`
	Identifier string `gorm:"column:account_identifier;size:160;index:idx_accounts_account_email"`
	Password   string `gorm:"column:account_password;size:256"`
	Status     string `gorm:"column:account_status;size:4"`
	Level      int    `gorm:"column:account_level"`
}

//||------------------------------------------------------------------------------------------------||
//|| Table Name
//||------------------------------------------------------------------------------------------------||

func (ModelAccount) TableName() string {
	return setup.Setup.Table
}
