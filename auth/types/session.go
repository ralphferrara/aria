package types

//||------------------------------------------------------------------------------------------------||
//|| Session Record
//||------------------------------------------------------------------------------------------------||

type SessionRecord struct {
	ID          int64
	Identifier  string
	Username    string
	Status      string
	Type        string
	Level       int
	Security    int
	Private     string
	PrivateHash string
	Public      string
	Created     int64
	Expires     int64
}
