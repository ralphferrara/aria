package db

//||------------------------------------------------------------------------------------------------||
//|| Status Codes
//||------------------------------------------------------------------------------------------------||

func UpdateStatusVerified(id uint) error {
	result := AuthDB().
		Model(&ModelAccount{}).
		Where("id_account = ?", id).
		Where("account_status = ?", "PEND").
		Update("account_status", "VERF")
	return result.Error
}

//||------------------------------------------------------------------------------------------------||
//|| Status Active
//||------------------------------------------------------------------------------------------------||

func UpdateStatusActive(id uint) error {
	result := AuthDB().
		Model(&ModelAccount{}).
		Where("id_account = ?", id).
		Where("account_status = ?", "VERF").
		Update("account_status", "ACTV")
	return result.Error
}
