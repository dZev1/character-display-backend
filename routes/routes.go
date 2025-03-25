package routes

import (
	characterHandlers "character-display-server/handlers/character_upload"
	loginHandlers "character-display-server/handlers/login"
	"character-display-server/middleware"
	"net/http"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", loginHandlers.Register)
	router.HandleFunc("POST /login", loginHandlers.Login)
	router.HandleFunc("GET /get_characters", characterHandlers.GetCharacters)
	router.HandleFunc("GET /get_all_characters", characterHandlers.GetAllCharacters)
	
	router.HandleFunc("POST /logout", middleware.Protected(loginHandlers.Logout))
	router.HandleFunc("POST /upload_character", middleware.Protected(characterHandlers.UploadCharacter))
	router.HandleFunc("GET /edit_character", middleware.Protected(characterHandlers.EditCharacter))
	router.HandleFunc("PUT /edit_character", middleware.Protected(characterHandlers.EditCharacter))
	router.HandleFunc("DELETE /delete_character", middleware.Protected(characterHandlers.DeleteCharacter))
	
	return router
}