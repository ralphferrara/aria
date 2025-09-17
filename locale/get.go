package locale

import (
	"fmt"
	"strings"
)

// ||------------------------------------------------------------------------------------------------||
// || GetTranslation
// ||------------------------------------------------------------------------------------------------||

func GetTranslation(section, key, lang string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	section = strings.ToUpper(section)
	key = strings.ToUpper(key)
	lang = strings.ToLower(lang)

	if sec, ok := translations[lang]; ok {
		if keys, ok := sec[section]; ok {
			if val, ok := keys[key]; ok {
				return val, nil
			}
		}
	}

	placeholder := fmt.Sprintf("--MISSINGTRANSLATION-%s-%s-%s", section, key, lang)
	return placeholder, fmt.Errorf("missing translation for %s.%s in %s", section, key, lang)
}
