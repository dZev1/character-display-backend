package main

import (
	"net/http"

	"github.com/dZev1/character-display/config"
	"github.com/dZev1/character-display/database"
	characterHandlers "github.com/dZev1/character-display/handlers/character_upload"
	loginHandlers "github.com/dZev1/character-display/handlers/login"
)

func main() {
	var err error
	router := http.NewServeMux()

	connStr, err := config.ReadConnStrEnv()
	if err != nil {
		panic(err)
	}

	err = database.InitDB(connStr)
	if err != nil {
		panic(err)
	}
	defer database.CloseDB()

	router.HandleFunc("POST /register", loginHandlers.Register)
	router.HandleFunc("POST /login", loginHandlers.Login)
	router.HandleFunc("POST /logout", loginHandlers.Logout)
	router.HandleFunc("POST /protected", loginHandlers.Protected)

	router.HandleFunc("POST /upload_character", characterHandlers.UploadCharacter)
	router.HandleFunc("GET /get_character", characterHandlers.GetCharacters)
	router.HandleFunc("GET /edit_character", characterHandlers.EditCharacter)
	router.HandleFunc("PUT /edit_character", characterHandlers.EditCharacter)
	
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err);
	}
}