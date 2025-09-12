package template

import (
	"fmt"
	"os"
)

//||------------------------------------------------------------------------------------------------||
//|| Template Markers
//||------------------------------------------------------------------------------------------------||

func Register(name string, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic("Template file not found: " + path)
	}
	Templates[name] = TemplateFile{
		Name: name,
		Path: path,
		Data: string(data),
	}
	fmt.Println("\033[32m[TEMP] - Registered template: " + name + " from " + path + "\033[0m")
}
