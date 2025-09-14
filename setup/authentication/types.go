package authentication

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

type AuthFields struct {
	ID         string
	Identifier string
	Username   string
	Password   string
	Level      string
	Status     string
}

//||------------------------------------------------------------------------------------------------||
//|| Auth: Globals
//||------------------------------------------------------------------------------------------------||

type Setup struct {
	Initialized bool
	Pepper      string
	CSRF        string
	Database    string
	Table       string
	Fields      AuthFields
}
