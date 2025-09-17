package template

import (
	"regexp"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| ParseMarkers: Extract translation markers from text
//||------------------------------------------------------------------------------------------------||

func ParseTranslations(input string) []ParsedTranslation {
	re := regexp.MustCompile(`\{\{::([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(input, -1)

	var markers []ParsedTranslation
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		raw := m[0]
		parts := strings.Split(m[1], ":")
		if len(parts) < 2 {
			continue
		}
		marker := ParsedTranslation{
			Raw:     raw,
			Section: strings.ToUpper(parts[0]),
			Key:     strings.ToUpper(parts[1]),
		}
		if len(parts) > 2 {
			marker.Casing = parts[2]
		}
		markers = append(markers, marker)
	}
	return markers
}
