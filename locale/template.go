package locale

import (
	"regexp"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| Parsed Translations
//||------------------------------------------------------------------------------------------------||

type ParsedTranslation struct {
	Raw     string
	Section string
	Key     string
	Casing  string
}

//||------------------------------------------------------------------------------------------------||
//|| ParseMarkers: Extract translation markers from text
//||------------------------------------------------------------------------------------------------||

func ParseTranslations(input string) []ParsedTranslation {
	re := regexp.MustCompile(`\{\{::([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(input, -1)

	var result []ParsedTranslation
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		raw := match[1]
		parts := strings.Split(raw, ":")
		if len(parts) < 2 {
			continue
		}

		marker := ParsedTranslation{
			Raw:     match[0], // full `{{::...}}`
			Section: parts[0],
			Key:     parts[1],
		}

		if len(parts) > 2 {
			marker.Casing = parts[2]
		}

		result = append(result, marker)
	}

	return result
}
