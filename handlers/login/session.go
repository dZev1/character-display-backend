package handlers

import (
	"errors"
	"net/http"

	"github.com/dZev1/character-display/database"
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

	csfrToken := r.Header.Get("X-CSRF-Token")
	if csfrToken != user.CSFRToken || csfrToken == "" {
		return AuthError
	}

	return nil
}
