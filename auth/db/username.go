package db

func CheckUsernameExists(username string) bool {
	var count int64
	err := AuthDB().
		Model(&ModelAccount{}).
		Where("account_username = ?", username).
		Count(&count).Error
	if err != nil {
		return true
	}
	return count > 0
}
