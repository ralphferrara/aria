package validate

import (
	"regexp"
	"strings"

	"github.com/ralphferrara/aria/app"
)

func IsEmailOrPhone(identifier string) string {
	// Trim spaces
	id := strings.TrimSpace(identifier)

	// Simple email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// Simple phone regex (digits, +, (), -, spaces allowed)
	phoneRegex := regexp.MustCompile(`^\+?[0-9\-\s\(\)]+$`)

	switch {
	case emailRegex.MatchString(id):
		return "email"
	case phoneRegex.MatchString(id):
		return "phone"
	default:
		return "unknown"
	}
}

func ValidatePhoneOrEmail(identifier string) error {
	it := IsEmailOrPhone(strings.TrimSpace(identifier))
	if it == "unknown" {
		return app.Err("Auth").Error("INVALID_IDENTIFIER")
	}
	if it == "email" && !IsValidEmail(identifier) {
		return app.Err("Auth").Error("INVALID_EMAIL")
	}
	if it == "phone" && !IsValidPhone(identifier) {
		return app.Err("Auth").Error("INVALID_PHONE")
	}
	return nil
}
