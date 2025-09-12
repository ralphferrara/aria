package locale

import "fmt"

//||------------------------------------------------------------------------------------------------||
//|| GetTranslation: Return a translation (json or txt)
//||------------------------------------------------------------------------------------------------||

func GetTranslation(section, term, lang string) (string, error) {
	if term == "" {
		if sectionData, ok := TextBlocks[section]; ok {
			if val, ok := sectionData[lang]; ok {
				return val, nil
			}
		}
		return "", fmt.Errorf("term is required")
	}

	if sectionMap, ok := Translations[section]; ok {
		key := fmt.Sprintf("%s.%s", lang, term)
		if val, ok := sectionMap[key]; ok {
			return val, nil
		}
	}

	return "", fmt.Errorf("translation not found")
}
