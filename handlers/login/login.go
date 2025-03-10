package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dZev1/character-display/database"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var err int
	if r.Method != http.MethodPost {
		err = http.StatusMethodNotAllowed
		http.Error(w, "invalid method", err)
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

	hashedPassword, _ := hashPassword(password)
	er := database.InsertUser(username, hashedPassword)
	if er != nil {
		http.Error(w, "could not insert user", http.StatusInternalServerError)
	}
	fmt.Fprintln(w, "user registered succesfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "invalid method", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	
	user, _ := database.GetUser(username)
	if ok, _ := database.IsInDatabase(username); !ok || checkPasswordHash(user.HashedPassword, password) {
		err := http.StatusUnauthorized
		http.Error(w, "user not found", err)
		return
	}

	sessionToken := generateToken(32)
	csfrToken := generateToken(32)
	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: sessionToken,
		Expires: time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name: "csfr_token",
		Value: csfrToken,
		HttpOnly: false,
	})

	user.SessionToken = sessionToken
	user.CSFRToken = csfrToken
	database.UpdateCookies(user)

	fmt.Fprintln(w, "login successful")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "unauthorized", er)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name: "csfr_token",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	username := r.FormValue("username")
	user, _ := database.GetUser(username)
	
	user.CSFRToken = ""
	user.SessionToken = ""
	
	database.UpdateCookies(user)

	fmt.Fprintln(w, "logged out succesfully.")
}

func Protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "invalid method", err)
		return
	}

	if err := Authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w,"unauthorized", er)
		return
	}

	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF validation succesful. Welcome, %s", username)
}