package main

import (
	"net/http"

	"github.com/dZev1/character-display/database"
	characterHandlers "github.com/dZev1/character-display/handlers/character_upload"
	loginHandlers "github.com/dZev1/character-display/handlers/login"
)

func main() {
	var err error

	connStr, err := connStrEnv()
	if err != nil {
		panic(err)
	}

	err = database.InitDB(connStr)
	if err != nil {
		panic(err)
	}
	defer database.CloseDB()

	http.HandleFunc("/register", loginHandlers.Register)
	http.HandleFunc("/login", loginHandlers.Login)
	http.HandleFunc("/logout", loginHandlers.Logout)
	http.HandleFunc("/protected", loginHandlers.Protected)

	http.HandleFunc("/upload_character", characterHandlers.UploadCharacter)
	
	http.ListenAndServe(":8080", nil)
}