package routes

import (
	"net/http"
	loginHandlers "character-display-server/handlers/login"
	characterHandlers "character-display-server/handlers/character_upload"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", loginHandlers.Register)
	router.HandleFunc("POST /login", loginHandlers.Login)
	router.HandleFunc("POST /logout", loginHandlers.Logout)
	router.HandleFunc("POST /protected", loginHandlers.Protected)

	router.HandleFunc("POST /upload_character", characterHandlers.UploadCharacter)
	router.HandleFunc("GET /get_characters", characterHandlers.GetCharacters)
	router.HandleFunc("GET /edit_character", characterHandlers.EditCharacter)
	router.HandleFunc("PUT /edit_character", characterHandlers.EditCharacter)
	router.HandleFunc("DELETE /delete_character", characterHandlers.DeleteCharacter)
	router.HandleFunc("GET /get_all_characters", characterHandlers.GetAllCharacters)

	return router
}