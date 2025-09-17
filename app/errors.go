package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ralphferrara/aria/locale"
)

//||------------------------------------------------------------------------------------------------||
//|| formatName
//||------------------------------------------------------------------------------------------------||

func formatName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ToUpper(name)
	return name
}

//||------------------------------------------------------------------------------------------------||
//|| Store
//||------------------------------------------------------------------------------------------------||

var ErrorsStore = ErrorsMaster{
	entries: make(map[string]ErrorsLibrary),
}

//||------------------------------------------------------------------------------------------------||
//|| Constants Master
//||------------------------------------------------------------------------------------------------||

type ErrorsMaster struct {
	entries map[string]ErrorsLibrary
}

//||------------------------------------------------------------------------------------------------||
//|| Constants Library
//||------------------------------------------------------------------------------------------------||

type ErrorsLibrary struct {
	name    string
	entries map[string]ErrorsEntry
}

//||------------------------------------------------------------------------------------------------||
//|| Constants Entry
//||------------------------------------------------------------------------------------------------||

type ErrorsEntry struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Fatal   bool   `json:"fatal"`
}

//||------------------------------------------------------------------------------------------------||
//|| Var
//||------------------------------------------------------------------------------------------------||

func Err(libraryName string) ErrorsLibrary {
	libraryName = formatName(libraryName)
	if !ErrorsStore.HasLibrary(libraryName) {
		ErrorsStore.AddLibrary(libraryName)
	}
	return ErrorsStore.GetLibrary(libraryName)
}

func ErrEntry(code string) ErrorsEntry {
	parts := strings.Split(code, ".")
	if len(parts) != 2 {
		fmt.Println("UNDEFINED ERROR MESSAGE:", code)
		return ErrorsEntry{
			Code:    "UNDEFINED_ERROR",
			Message: "UNDEFINED ERROR MESSAGE",
			Fatal:   false,
		}
	}
	return Err(parts[0]).Get(parts[1])
}

func ErrMessage(code string) string {
	parts := strings.Split(code, ".")
	if len(parts) != 2 {
		fmt.Println("UNDEFINED ERROR MESSAGE:", code)
		return "UNDEFINED_ERROR"
	}
	return Err(parts[0]).Message(parts[1])
}

func ErrFatal(code string) bool {
	parts := strings.Split(code, ".")
	if len(parts) != 2 {
		return false
	}
	return Err(parts[0]).IsFatal(parts[1])
}

//||------------------------------------------------------------------------------------------------||
//|| Get
//||------------------------------------------------------------------------------------------------||

func (lib ErrorsLibrary) Get(code string) ErrorsEntry {
	code = formatName(code)
	if entry, exists := lib.entries[code]; exists {
		return entry
	}
	fmt.Println("UNDEFINED ERROR MESSAGE:", code)
	return ErrorsEntry{
		Code:    "UNDEFINED_ERROR",
		Message: "UNDEFINED ERROR MESSAGE",
		Fatal:   false,
	}
}

func (lib ErrorsLibrary) Code(code string) string {
	code = formatName(code)
	fmt.Println("Error:", code)
	if entry, exists := lib.entries[code]; exists {
		return fmt.Sprintf("%s.%s", strings.ToUpper(lib.name), entry.Code)
	}
	fmt.Println("UNDEFINED ERROR CODE:", lib.name, code)
	return "ARIA.UNDEFINED_ERROR"
}

func (lib ErrorsLibrary) Message(code string) string {
	code = formatName(code)
	if entry, exists := lib.entries[code]; exists {
		return entry.Message

	}
	fmt.Println("UNDEFINED ERROR MESSAGE:", code)
	return "UNDEFINED ERROR MESSAGE"
}

func (lib ErrorsLibrary) Locale(code string, lang string) string {
	code = formatName(code)
	if entry, exists := lib.entries[code]; exists {
		translate, err := locale.GetTranslation("errors", entry.Code, lang)
		if err == nil {
			return translate
		}
	}
	return "UNDEFINED LOCALE ERROR MESSAGE"
}

func (lib ErrorsLibrary) Error(code string) error {
	code = formatName(code)
	if entry, exists := lib.entries[code]; exists {
		return fmt.Errorf("%s.%s", strings.ToUpper(lib.name), entry.Code)
	}
	return errors.New("UNDEFINED_ERROR")
}

func (lib ErrorsLibrary) IsFatal(code string) bool {
	code = formatName(code)
	if entry, exists := lib.entries[code]; exists {
		return entry.Fatal
	}
	return false
}

func (lib ErrorsLibrary) List() []ErrorsEntry {
	list := []ErrorsEntry{}
	for _, entry := range lib.entries {
		list = append(list, entry)
	}
	return list
}

//||------------------------------------------------------------------------------------------------||
//|| Var
//||------------------------------------------------------------------------------------------------||

func (lib ErrorsLibrary) Add(code, message string, fatal bool) {
	code = formatName(code)
	lib.entries[code] = ErrorsEntry{
		Code:    code,
		Message: message,
		Fatal:   fatal,
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Errors Master
//||------------------------------------------------------------------------------------------------||

func (m ErrorsMaster) AddLibrary(libName string) {
	libName = formatName(libName)
	if _, exists := m.entries[libName]; !exists {
		m.entries[libName] = ErrorsLibrary{name: libName, entries: make(map[string]ErrorsEntry)}
	}
}

func (m ErrorsMaster) HasLibrary(code string) bool {
	code = formatName(code)
	_, ok := m.entries[code]
	return ok
}

func (m ErrorsMaster) GetLibrary(libName string) ErrorsLibrary {
	libName = formatName(libName)
	if _, exists := m.entries[libName]; !exists {
		m.entries[libName] = ErrorsLibrary{name: libName, entries: make(map[string]ErrorsEntry)}
	}
	return m.entries[libName]
}
