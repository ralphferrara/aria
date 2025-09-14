package db

func DeleteAccount(id int64) error {
	result := AuthDB().Delete(&ModelAccount{}, id)
	return result.Error
}
