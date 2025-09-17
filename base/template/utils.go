package template

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"net/http"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| Replace Marker in String
//||------------------------------------------------------------------------------------------------||

func replaceMarker(data string, marker string, value string) string {
	placeholder := fmt.Sprintf("{{%s}}", marker)
	return strings.ReplaceAll(data, placeholder, value)
}

//||------------------------------------------------------------------------------------------------||
//|| Helper
//||------------------------------------------------------------------------------------------------||

func ReplaceMarker(data string, marker string, value string) string {
	return replaceMarker(data, marker, value)
}

//||------------------------------------------------------------------------------------------------||
//|| Serves the HTML file with dynamic replacements
//||------------------------------------------------------------------------------------------------||

func GetLanguageFromRequest(r *http.Request) string {
	lang := "en"
	if al := r.Header.Get("Accept-Language"); al != "" {
		primary := strings.SplitN(al, ",", 2)[0]
		if code := strings.SplitN(primary, "-", 2)[0]; code != "" {
			lang = strings.ToLower(code)
		}
	}
	return lang
}
