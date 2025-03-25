package middleware

import (
	"errors"
	"net/http"

	"character-display-server/database"
)

var AuthError = errors.New("Unauthorized")


func Authorize(r *http.Request) error {
	username := r.FormValue("username")
	user, ok := database.GetUser(username)
	if ok != nil {
		return AuthError
	}
	sessionToken, err := r.Cookie("session_token")
	if err != nil || sessionToken.Value == "" || sessionToken.Value != user.SessionToken {
		return AuthError
	}

	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken != user.CSRFToken || csrfToken == "" {
		return AuthError
	}

	return nil
}