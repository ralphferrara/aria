package locale

//||------------------------------------------------------------------------------------------------||
//|| Globals
//||------------------------------------------------------------------------------------------------||

var (
	Supported = []string{"en", "es", "fr", "de", "it", "pt", "zh"}
	Directory string
)

//||------------------------------------------------------------------------------------------------||
//|| Init: Creates LocaleWrapper and loads translations
//||------------------------------------------------------------------------------------------------||

func Init(dir string) LocaleWrapper {
	Directory = dir

	if err := LoadRendered(dir); err != nil {
		panic("failed to load translations: " + err.Error())
	}

	return LocaleWrapper{
		Directory: dir,
		Get: func(section, term, lang string) (string, error) {
			return GetTranslation(section, term, lang)
		},
		Load: LoadRendered,
	}
}
