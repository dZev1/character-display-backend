package handlers

import (
	"fmt"
	"net/http"
	"time"

	"character-display-server/database"
	"character-display-server/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var err int

	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not process form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) < 3 || len(password) < 8 {
		err = http.StatusNotAcceptable
		http.Error(w, "invalid username/password", err)
		return
	}

	if ok, _ := database.IsInDatabase(username); ok {
		err := http.StatusConflict
		http.Error(w, "user already exists", err)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	er := database.InsertUser(username, hashedPassword)
	if er != nil {
		http.Error(w, "could not insert user", http.StatusInternalServerError)
	}
	fmt.Fprintln(w, "user registered succesfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, _ := database.GetUser(username)
	if ok, _ := database.IsInDatabase(username); !ok || utils.CheckPasswordHash(user.HashedPassword, password) {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		HttpOnly: false,
	})

	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken

	err := database.UpdateCookies(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "login successful")
}

func Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	username := r.FormValue("username")
	user, _ := database.GetUser(username)

	user.CSRFToken = ""
	user.SessionToken = ""

	err := database.UpdateCookies(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "logged out succesfully.")
}
