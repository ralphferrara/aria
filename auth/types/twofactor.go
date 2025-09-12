package types

import (
	"time"
)

//||------------------------------------------------------------------------------------------------||
//|| Two Factor Verification Email
//||------------------------------------------------------------------------------------------------||

type TwoFactorVerification struct {
	Code       string    `json:"code"`
	Key        string    `json:"key"`
	Type       string    `json:"type"`
	Identifier string    `json:"identifier"`
	Attempts   int       `json:"attempts"`
	Created    time.Time `json:"created"`
	Expires    time.Time `json:"expires"`
}
