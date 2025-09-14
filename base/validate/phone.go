package validate

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

//||------------------------------------------------------------------------------------------------||
//|| IsValidPhone
//||------------------------------------------------------------------------------------------------||

func IsValidPhone(number string) bool {
	number = strings.TrimSpace(number)

	region := ""
	if !strings.HasPrefix(number, "+") {
		region = "US"
	}

	parsed, err := phonenumbers.Parse(number, region)
	if err != nil {
		return false
	}

	return phonenumbers.IsValidNumber(parsed)
}

//||------------------------------------------------------------------------------------------------||
//|| FormatPhone
//||------------------------------------------------------------------------------------------------||

func FormatPhone(number string) string {
	number = strings.TrimSpace(number)

	region := ""
	if !strings.HasPrefix(number, "+") {
		region = "US"
	}

	parsed, err := phonenumbers.Parse(number, region)
	if err != nil {
		return ""
	}

	if !phonenumbers.IsValidNumber(parsed) {
		return ""
	}

	return phonenumbers.Format(parsed, phonenumbers.E164)
}
