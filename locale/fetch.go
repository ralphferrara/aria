package locale

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"net/http"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| Requested Language
//||------------------------------------------------------------------------------------------------||

func Request(r *http.Request) string {
	langHeader := r.Header.Get("Accept-Language")
	if langHeader == "" {
		return "en" // fallback to English
	}

	// Example: "en-US,en;q=0.9,es;q=0.8"
	languages := strings.Split(langHeader, ",")
	if len(languages) == 0 {
		return "en"
	}

	// Get the first preferred language
	return strings.TrimSpace(strings.Split(languages[0], ";")[0])
}
