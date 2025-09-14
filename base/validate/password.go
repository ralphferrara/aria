package validate

import "github.com/ralphferrara/aria/app"

func IsValidPassword(password string) error {
	if password == "" || len(password) < 8 {
		return app.Err("Validation").Error("PASSWORD_TOO_SHORT")
	}
	return nil
}
