package types

//||------------------------------------------------------------------------------------------------||
//|| Session Record
//||------------------------------------------------------------------------------------------------||

type AuthMeRecord struct {
	ID         int64
	Identifier string
	Username   string
	Status     string
	Type       string
	Level      int
	Created    int64
	Expires    int64
}
