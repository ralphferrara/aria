package auth

import "github.com/ralphferrara/aria/app"

//||------------------------------------------------------------------------------------------------||
//|| Initialize Auth Constants
//||------------------------------------------------------------------------------------------------||

func InitErrors() {
	app.Err("Auth").Add("ACCOUNT_ALREADY_CREATED", "Account is already created", false)
	app.Err("Auth").Add("ACCOUNT_NOT_PENDING", "Account is already created", false)
	app.Err("Auth").Add("ACCOUNT_LOOKUP_FAILED", "Could not re-fetch account after update", false)
	app.Err("Auth").Add("ACCOUNT_DELETE_FAILED", "Failed to delete account", false)
	app.Err("Auth").Add("ACCOUNT_CREATE_FAILED", "Failed to create account", false)
	app.Err("Auth").Add("ACCOUNT_NOT_FOUND", "Account Not Found", false)
	app.Err("Auth").Add("ACCOUNT_TOKEN_MISMATCH", "Token does not match", false)

	app.Err("Auth").Add("USERNAME_GEN_FAILED", "Failed to generate username", false)

	app.Err("Auth").Add("PASSWORD_TOO_SHORT", "Password must be at least 8 characters long", false)
	app.Err("Auth").Add("PASSWORD_GEN_FAILED", "Could not generate password", false)
	app.Err("Auth").Add("PASSWORD_UPDATE_FAILED", "Password Update Failed", false)
	app.Err("Auth").Add("PASSWORD_MISMATCH", "Passwords do not match", false)

	app.Err("Auth").Add("SESSION_GEN_FAILED", "Failed to create session", false)
	app.Err("Auth").Add("SESSION_LOOKUP_FAILED", "Could not retrieve session", false)
	app.Err("Auth").Add("MISSING_SESSION_COOKIE", "Missing session cookie", false)

	app.Err("Auth").Add("PRIVPUB_FAILED", "Failed to create an encryption key pair", false)
	app.Err("Auth").Add("BAD_BIP39", "Invalid BIP39 Keywotd", false)
	app.Err("Auth").Add("BIP39_GEN_FAILED", "Failed to generate BIP39 keyword list", false)
	app.Err("Auth").Add("PRIVPUB_MISMATCH", "Private/Public Keypair do not match", false)
	app.Err("Auth").Add("PRIVPUB_CHECKEY_FAILED", "Failed to validate Private/Public Keypair", false)

	app.Err("Auth").Add("INVALID_EMAIL", "Email address is not valid", false)
	app.Err("Auth").Add("INVALID_ACCOUNT_TYPE", "Account type is not valid", false)
	app.Err("Auth").Add("INVALID_PHONE", "Phone number is not valid", false)
	app.Err("Auth").Add("INVALID_IDENTIFIER", "Identifier is not a valid phone number or email address", false)
	app.Err("Auth").Add("INVALID_CREDENTIALS", "Account Not Found", false)

	app.Err("Auth").Add("TF_MISSING_CODE", "Invalid Code/Token", false)
	app.Err("Auth").Add("TF_INVALID_TOKEN", "Invalid or expired token", false)
	app.Err("Auth").Add("TF_INVALID_RECORD", "Invalid stored record", false)
	app.Err("Auth").Add("TF_TOO_MANY_ATTEMPTS", "Too many attempts", false)
	app.Err("Auth").Add("TF_CODE_MISMATCH", "Invalid code", false)
	app.Err("Auth").Add("TF_TOKEN_EXPIRED", "Token expired", false)
}
