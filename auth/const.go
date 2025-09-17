package auth

import "github.com/ralphferrara/aria/app"

//||------------------------------------------------------------------------------------------------||
//|| Initialize Auth Constants
//||------------------------------------------------------------------------------------------------||

func InitConstants() {
	//||------------------------------------------------------------------------------------------------||
	//|| Account Statuses
	//||------------------------------------------------------------------------------------------------||
	app.Constants("AccountStatus").AddCode("Pending", "PEND", "Account has not been setup.")
	app.Constants("AccountStatus").AddCode("Verified", "VERF", "Account identifier is verified.")
	app.Constants("AccountStatus").AddCode("Active", "ACTV", "Account is active.")
	app.Constants("AccountStatus").AddCode("Suspended", "SUSP", "Account is suspended.")
	app.Constants("AccountStatus").AddCode("Deleted", "DELD", "Account is deleted.")
	//||------------------------------------------------------------------------------------------------||
	//|| Two Factor Types
	//||------------------------------------------------------------------------------------------------||
	app.Constants("TwoFactorType").AddCode("Reset", "RESET", "Verify For Password Reset")
	app.Constants("TwoFactorType").AddCode("Account", "ACCOUNT", "Verify New Account")

}
