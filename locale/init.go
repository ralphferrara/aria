package locale

//||------------------------------------------------------------------------------------------------||
//|| Globals
//||------------------------------------------------------------------------------------------------||

var (
	Directory    = "./locales"
	Translations = make(map[string]map[string]string)
	TextBlocks   = make(map[string]map[string]string)
)

//||------------------------------------------------------------------------------------------------||
//|| Init: Creates LocaleWrapper and returns helper
//||------------------------------------------------------------------------------------------------||

func Init(dir string) LocaleWrapper {
	Directory = dir
	local := LocaleWrapper{
		Directory:    dir,
		Translations: Translations,
		TextBlocks:   TextBlocks,
		Get: func(section, term, lang string) (string, error) {
			return GetTranslation(section, term, lang)
		},
		Load: Load,
	}
	local.Load()
	return local
}
