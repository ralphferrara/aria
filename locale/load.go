package locale

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| Load: Read translated files from /translated
//||------------------------------------------------------------------------------------------------||

func Load() {
	files, err := os.ReadDir(filepath.Join(Directory, "translated"))
	if err != nil {
		return
	}

	for _, f := range files {
		parts := strings.SplitN(f.Name(), ".", 3)
		if len(parts) != 3 {
			continue
		}

		lang := parts[0]
		section := parts[1]
		ext := parts[2]
		path := filepath.Join(Directory, "translated", f.Name())

		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		switch ext {
		case "txt":
			if _, ok := TextBlocks[section]; !ok {
				TextBlocks[section] = make(map[string]string)
			}
			TextBlocks[section][lang] = string(content)

		case "json":
			var data map[string]string
			if err := json.Unmarshal(content, &data); err != nil {
				continue
			}

			if _, ok := Translations[section]; !ok {
				Translations[section] = make(map[string]string)
			}

			for key, val := range data {
				composite := fmt.Sprintf("%s.%s", lang, key)
				Translations[section][composite] = val
			}
		}
	}
}
