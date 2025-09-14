package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ralphferrara/aria/app"
	"github.com/ralphferrara/aria/auth/db"
	"github.com/ralphferrara/aria/auth/types"
	"github.com/ralphferrara/aria/base/random"
)

//||------------------------------------------------------------------------------------------------||
//|| Get the Account Record
//||------------------------------------------------------------------------------------------------||

func SessionCreate(identifier string, account *db.ModelAccount) (string, error) {

	//||------------------------------------------------------------------------------------------------||
	//|| Generate a Random Token
	//||------------------------------------------------------------------------------------------------||

	sessionToken := random.UUIDString()

	//||------------------------------------------------------------------------------------------------||
	//|| Check if Exists
	//||------------------------------------------------------------------------------------------------||

	session := types.SessionRecord{
		ID:         account.ID,
		Identifier: identifier,
		Username:   account.Username,
		Status:     account.Status,
		Level:      account.Level,
		Created:    time.Now().Unix(),
		Expires:    time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Marshal
	//||------------------------------------------------------------------------------------------------||

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		fmt.Println("[Session] Failed to marshal session:", err)
		return "", err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Save to Redis
	//||------------------------------------------------------------------------------------------------||

	err = app.CacheRedis["auth"].Set("session::"+sessionToken, sessionJSON, 30*24*time.Hour)
	if err != nil {
		fmt.Println("[Session] Failed to save session to Redis:", err)
		return "", err
	}
	//||------------------------------------------------------------------------------------------------||
	//|| Done
	//||------------------------------------------------------------------------------------------------||

	fmt.Println("[Session] Token set for account:", sessionToken)
	return sessionToken, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Fetch Session
//||------------------------------------------------------------------------------------------------||

func FetchSession(sessionID string) (types.SessionRecord, error) {
	//||------------------------------------------------------------------------------------------------||
	//|| Get the Session from the Database
	//||------------------------------------------------------------------------------------------------||

	sessionJSON, err := app.CacheRedis["auth"].Get("session::" + sessionID)
	if err != nil {
		return types.SessionRecord{}, fmt.Errorf("failed to fetch session: %w", err)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Return the Session
	//||------------------------------------------------------------------------------------------------||

	var session types.SessionRecord
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return types.SessionRecord{}, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Check Session
	//||------------------------------------------------------------------------------------------------||

	if time.Unix(session.Expires, 0).Before(time.Now()) {
		return types.SessionRecord{}, fmt.Errorf("session expired")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Return the Session
	//||------------------------------------------------------------------------------------------------||

	return session, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Update the Session
//||------------------------------------------------------------------------------------------------||

func UpdateSession(sessionToken string, session types.SessionRecord) error {

	//||------------------------------------------------------------------------------------------------||
	//|| Marshal Session Data
	//||------------------------------------------------------------------------------------------------||

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		fmt.Println("[Session] Failed to marshal session:", err)
		return err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Save to Redis (overwrite)
	//||------------------------------------------------------------------------------------------------||

	err = app.CacheRedis["auth"].Set("session::"+sessionToken, sessionJSON, 30*24*time.Hour)
	if err != nil {
		fmt.Println("[Session] Failed to update session in Redis:", err)
		return err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Done
	//||------------------------------------------------------------------------------------------------||

	fmt.Println("[Session] Updated session for account:", sessionToken)
	return nil
}

//||------------------------------------------------------------------------------------------------||
//|| Create the
//||------------------------------------------------------------------------------------------------||

func WriteSessionCookie(w http.ResponseWriter, sessionToken string) {
	//||------------------------------------------------------------------------------------------------||
	//|| Set the Cookie
	//||------------------------------------------------------------------------------------------------||

	if sessionToken == "" {
		fmt.Println("[Session] No session token provided")
		return
	}

	fmt.Println("[Session] Setting cookie with token:", sessionToken)

	//||------------------------------------------------------------------------------------------------||
	//|| UI Cookie
	//||------------------------------------------------------------------------------------------------||

	http.SetCookie(w, &http.Cookie{
		Name:     "session_ui",
		Value:    "1",
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 30,
		HttpOnly: false,
	})

	//||------------------------------------------------------------------------------------------------||
	//|| Create and set the cookie
	//||------------------------------------------------------------------------------------------------||

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 30,
	})

}

//||------------------------------------------------------------------------------------------------||
//|| Delete Session
//||------------------------------------------------------------------------------------------------||

func DeleteSession(sessionToken string) error {
	return app.CacheRedis["auth"].Del("session::" + sessionToken)
}

//||------------------------------------------------------------------------------------------------||
//|| Clear Session Cookie
//||------------------------------------------------------------------------------------------------||

func ClearSessionCookie(w http.ResponseWriter) {

	// Helper to generate an expired cookie with matching attributes
	expiredCookie := func(name string, httpOnly bool) *http.Cookie {
		return &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",             // Must match original Path
			Domain:   "",              // Set if you explicitly set Domain
			Expires:  time.Unix(0, 0), // Past date
			MaxAge:   -1,              // Delete immediately
			HttpOnly: httpOnly,
			Secure:   true, // Must match original Secure flag
			SameSite: http.SameSiteLaxMode,
		}
	}

	// Expire both cookies
	http.SetCookie(w, expiredCookie("session", true))
	http.SetCookie(w, expiredCookie("session_ui", false))

	fmt.Println("[Session] Cleared session and session_ui cookies")
}
