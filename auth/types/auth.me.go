package types

import "time"

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
	Created    time.Time
	LastLogin  time.Time
}
