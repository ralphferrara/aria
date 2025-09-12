package convert

import (
	"fmt"
	"time"
)

//||------------------------------------------------------------------------------------------------||
//|| Type Names
//||------------------------------------------------------------------------------------------------||

var monthNames = map[string][]string{
	"en": {
		"", // placeholder for 0
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	},
	"fr": {
		"",
		"janvier", "février", "mars", "avril", "mai", "juin",
		"juillet", "août", "septembre", "octobre", "novembre", "décembre",
	},
	"es": {
		"",
		"enero", "febrero", "marzo", "abril", "mayo", "junio",
		"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
	},
}

//||------------------------------------------------------------------------------------------------||
//|| Format
//||------------------------------------------------------------------------------------------------||

func FormatMonthYear(t time.Time, lang string) string {
	m := int(t.Month())
	y := t.Year()
	names, ok := monthNames[lang]
	if !ok {
		names = monthNames["en"]
	}
	// same order for all langs: Month then Year
	return fmt.Sprintf("%s %d", names[m], y)
}
