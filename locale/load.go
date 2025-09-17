package locale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

// ||------------------------------------------------------------------------------------------------||
// || Translation Store
// ||------------------------------------------------------------------------------------------------||
var (
	translations = make(map[string]map[string]map[string]string) // lang -> section -> key -> value
	mu           sync.RWMutex
)

// ||------------------------------------------------------------------------------------------------||
// || LoadRendered: Scan `<dir>/.rendered` and cache translations
// ||------------------------------------------------------------------------------------------------||
func LoadRendered(dir string) error {
	// Always look inside the `.rendered` subdirectory
	renderedDir := filepath.Join(dir, ".rendered")

	files, err := filepath.Glob(filepath.Join(renderedDir, "*.*"))
	if err != nil {
		return fmt.Errorf("failed to list rendered files: %w", err)
	}

	var countJSON = 0
	var countTXT = 0

	for _, file := range files {
		base := filepath.Base(file)
		parts := strings.SplitN(base, ".", 3) // lang.name.ext
		if len(parts) != 3 {
			continue
		}

		lang := strings.ToLower(parts[0])
		name := strings.ToUpper(parts[1])
		ext := strings.ToLower(parts[2])

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		mu.Lock()
		if _, ok := translations[lang]; !ok {
			translations[lang] = make(map[string]map[string]string)
		}

		switch ext {
		case "json":
			section := name
			if _, ok := translations[lang][section]; !ok {
				translations[lang][section] = make(map[string]string)
			}

			var parsed map[string]string
			if err := json.Unmarshal(data, &parsed); err != nil {
				mu.Unlock()
				return fmt.Errorf("failed to parse JSON in %s: %w", file, err)
			}
			for k, v := range parsed {
				key := strings.ToUpper(k)
				translations[lang][section][key] = v
				countJSON++
			}

		case "txt":
			section := "BLURB"
			if _, ok := translations[lang][section]; !ok {
				translations[lang][section] = make(map[string]string)
			}
			translations[lang][section][name] = strings.TrimSpace(string(data))
			countTXT++

		default:
			// ignore unsupported ext
		}
		mu.Unlock()
	}
	fmt.Printf("[locale] loaded %d entries\n", countJSON+countTXT)
	return nil
}
