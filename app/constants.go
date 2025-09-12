package app

import "github.com/ralphferrara/aria/locale"

//||------------------------------------------------------------------------------------------------||
//|| Store
//||------------------------------------------------------------------------------------------------||

var ConstantsStore = ConstantsMaster{
	entries: make(map[string]ConstantsLibrary),
}

//||------------------------------------------------------------------------------------------------||
//|| Constants Master
//||------------------------------------------------------------------------------------------------||

type ConstantsMaster struct {
	entries map[string]ConstantsLibrary
}

//||------------------------------------------------------------------------------------------------||
//|| Constants Library
//||------------------------------------------------------------------------------------------------||

type ConstantsLibrary struct {
	entries map[string]ConstantsEntry
}

//||------------------------------------------------------------------------------------------------||
//|| Constants Entry
//||------------------------------------------------------------------------------------------------||

type ConstantsEntry struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	ValueString string `json:"value"`
	ValueInt    int    `json:"value_int"`
	ValueBool   bool   `json:"value_bool"`
}

//||------------------------------------------------------------------------------------------------||
//|| Var
//||------------------------------------------------------------------------------------------------||

func Constants(libraryName string) ConstantsLibrary {
	if !ConstantsStore.HasLibrary(libraryName) {
		ConstantsStore.AddLibrary(libraryName)
	}
	return ConstantsStore.GetLibrary(libraryName)
}

//||------------------------------------------------------------------------------------------------||
//|| Get
//||------------------------------------------------------------------------------------------------||

func (lib ConstantsLibrary) Get(name string) ConstantsEntry {
	name = formatName(name)
	if entry, exists := lib.entries[name]; exists {
		return entry
	}
	return ConstantsEntry{
		Type:        "undefined",
		Name:        name,
		Code:        "UNDEFINED",
		ValueString: "UNDEFINED CONSTANT",
	}
}

func (lib ConstantsLibrary) Bool(name string) bool {
	name = formatName(name)
	if entry, exists := lib.entries[name]; exists {
		return entry.ValueBool
	}
	return false
}

func (lib ConstantsLibrary) Int(name string) int {
	name = formatName(name)
	if entry, exists := lib.entries[name]; exists {
		return entry.ValueInt
	}
	return 0
}

func (lib ConstantsLibrary) String(name string) string {
	name = formatName(name)
	if entry, exists := lib.entries[name]; exists {
		return entry.ValueString
	}
	return ""
}

func (lib ConstantsLibrary) Locale(code string, lang string) string {
	code = formatName(code)
	if entry, exists := lib.entries[code]; exists {
		translate, err := locale.GetTranslation("constants", entry.Code, lang)
		if err == nil {
			return translate
		}
	}
	return "UNDEFINED LOCALE ERROR MESSAGE"
}

func (lib ConstantsLibrary) Code(name string) string {
	name = formatName(name)
	if entry, exists := lib.entries[name]; exists {
		return entry.Code
	}
	return ""
}

func (lib ConstantsLibrary) Desciption(name string) string {
	name = formatName(name)
	return lib.String(name)
}

func (lib ConstantsLibrary) List() []ConstantsEntry {
	list := []ConstantsEntry{}
	for _, entry := range lib.entries {
		list = append(list, entry)
	}
	return list
}

//||------------------------------------------------------------------------------------------------||
//|| Var
//||------------------------------------------------------------------------------------------------||

func (lib ConstantsLibrary) AddString(name, value string) {
	name = formatName(name)
	lib.entries[name] = ConstantsEntry{
		Type:        "string",
		Name:        name,
		ValueString: value,
	}
}

func (lib ConstantsLibrary) AddCode(name, code, description string) {
	name = formatName(name)
	lib.entries[name] = ConstantsEntry{
		Type:        "string",
		Name:        name,
		Code:        code,
		ValueString: description,
	}
}

func (lib ConstantsLibrary) AddInt(name string, value int) {
	name = formatName(name)
	lib.entries[name] = ConstantsEntry{
		Type:     "int",
		Name:     name,
		ValueInt: value,
	}
}

func (lib ConstantsLibrary) AddBool(name string, value bool) {
	name = formatName(name)
	lib.entries[name] = ConstantsEntry{
		Type:      "bool",
		Name:      name,
		ValueBool: value,
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Constants Master
//||------------------------------------------------------------------------------------------------||

func (m ConstantsMaster) AddLibrary(name string) {
	name = formatName(name)
	if _, exists := m.entries[name]; !exists {
		m.entries[name] = ConstantsLibrary{entries: make(map[string]ConstantsEntry)}
	}
}

func (m ConstantsMaster) HasLibrary(name string) bool {
	name = formatName(name)
	_, ok := m.entries[name]
	return ok
}

func (m ConstantsMaster) GetLibrary(name string) ConstantsLibrary {
	name = formatName(name)
	if _, exists := m.entries[name]; !exists {
		m.entries[name] = ConstantsLibrary{entries: make(map[string]ConstantsEntry)}
	}
	return m.entries[name]
}
