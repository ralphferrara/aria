package validate

//||------------------------------------------------------------------------------------------------||
//|| Validate Date (YYYY, MM, DD)
//||------------------------------------------------------------------------------------------------||

func IsValidDate(year, month, day int) bool {
	// year must be positive (4-digit typical, but allow >0)
	if year <= 0 {
		return false
	}
	// month must be 1..12
	if month < 1 || month > 12 {
		return false
	}
	// day must be >=1
	if day < 1 {
		return false
	}

	// days per month (default February = 28)
	daysInMonth := [...]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// leap year adjustment for February
	if (year%4 == 0 && year%100 != 0) || (year%400 == 0) {
		daysInMonth[1] = 29
	}

	// day must not exceed month length
	if day > daysInMonth[month-1] {
		return false
	}

	return true
}
