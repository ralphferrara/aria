package locale

//||------------------------------------------------------------------------------------------------||
//|| Types
//||------------------------------------------------------------------------------------------------||

type LocaleWrapper struct {
	Directory string
	Get       func(section, term, lang string) (string, error)
	Load      func(dir string) error
}
