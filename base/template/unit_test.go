package template

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"os"
	"testing"

	"github.com/ralphferrara/aria/locale"
)

//||------------------------------------------------------------------------------------------------||
//|| Unit Test: Compile with Translation + Marker Replacement
//||------------------------------------------------------------------------------------------------||

func TestCompileTemplate(t *testing.T) {

	//||------------------------------------------------------------------------------------------------||
	//|| Setup: Register a template with translations
	//||------------------------------------------------------------------------------------------------||

	testFile := "test_template.txt"
	testContent := "Hello {{phrases.name}}, welcome to {{phrases.site}}! {{STATIC}}"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test template file: %v", err)
	}
	defer os.Remove(testFile)

	Register("test", testFile)

	//||------------------------------------------------------------------------------------------------||
	//|| Mock Translation Data
	//||------------------------------------------------------------------------------------------------||

	locale.Translations = map[string]map[string]string{
		"phrases": {
			"en.name": "John",
			"en.site": "ComplyAge",
		},
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create Instance & Compile
	//||------------------------------------------------------------------------------------------------||

	tmpl := Create("test", "en")
	tmpl.Add("STATIC", "Enjoy your stay.")

	result := tmpl.Compile()

	expected := "Hello John, welcome to ComplyAge! Enjoy your stay."
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}
