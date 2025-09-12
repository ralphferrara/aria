package locale

//||------------------------------------------------------------------------------------------------||
//|| Types
//||------------------------------------------------------------------------------------------------||

type LocaleWrapper struct {
	Directory    string
	Translations map[string]map[string]string
	TextBlocks   map[string]map[string]string
	Get          func(section, term, lang string) (string, error)
	Load         func()
}
