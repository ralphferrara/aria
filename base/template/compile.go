package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ralphferrara/aria/locale"
)

//||------------------------------------------------------------------------------------------------||
//|| Template Markers
//||------------------------------------------------------------------------------------------------||

func (t *TemplateInstance) Compile() string {
	// Handle translation markers
	for _, marker := range ParseTranslations(t.Data) {
		value, err := locale.GetTranslation(marker.Section, marker.Key, t.Language)
		if err != nil {
			fmt.Printf("Missing translation for %s:%s (%s): %v\n",
				marker.Section, marker.Key, t.Language, err)
			continue
		}

		switch strings.ToUpper(marker.Casing) {
		case "UPPER":
			value = strings.ToUpper(value)
		case "LOWER":
			value = strings.ToLower(value)
		case "TITLE":
			value = strings.Title(value)
		}

		t.Data = strings.ReplaceAll(t.Data, marker.Raw, value)
	}

	// Handle normal markers (non-translation)
	re := regexp.MustCompile(`\{\{([^}:]+)\}\}`) // matches {{NAME}} but not {{::...}}
	matches := re.FindAllStringSubmatch(t.Data, -1)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		raw := m[0]
		name := m[1]

		for _, marker := range t.Markers {
			if marker.Marker == name {
				t.Data = strings.ReplaceAll(t.Data, raw, marker.Value)
				break
			}
		}
	}

	return t.Data
}
