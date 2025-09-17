package template

import (
	"fmt"
	"os"
)

// ||------------------------------------------------------------------------------------------------||
// || Create: Load template instance, disk in dev, memory in production
// ||------------------------------------------------------------------------------------------------||

func Create(name string) TemplateInstance {
	env := os.Getenv("ENV_MODE")

	// Production → use cached Templates
	if env == "production" {
		if tmpl, ok := Templates[name]; ok {
			return TemplateInstance{
				Data:     tmpl.Data,
				Markers:  []TemplateMarker{},
				Language: "en",
			}
		}
		fmt.Printf("[template] missing template in memory: %s\n", name)
		return TemplateInstance{Data: "", Markers: []TemplateMarker{}, Language: "en"}
	}

	// Dev → reload from disk using registered path if possible
	if tmpl, ok := Templates[name]; ok {
		data, err := os.ReadFile(tmpl.Path)
		if err != nil {
			fmt.Printf("[template] failed to read %s: %v\n", tmpl.Path, err)
			return TemplateInstance{Data: "", Markers: []TemplateMarker{}, Language: "en"}
		}
		return TemplateInstance{
			Data:     string(data),
			Markers:  []TemplateMarker{},
			Language: "en",
		}
	}

	fmt.Printf("[template] template not registered: %s\n", name)
	return TemplateInstance{Data: "", Markers: []TemplateMarker{}, Language: "en"}
}
