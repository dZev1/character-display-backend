package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dZev1/character-display/database"
	loginHandlers "github.com/dZev1/character-display/handlers/login"
	"github.com/dZev1/character-display/models"
)

func UploadCharacter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "method not allowed", err)
		return
	}

	err := r.ParseMultipartForm(10<<20)
	if err != nil {
		er := http.StatusBadRequest
		http.Error(w, err.Error(), er)
		return
	}

	if err := loginHandlers.Authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w,"unauthorized", er)
		return
	}

	charJSON := r.FormValue("char_json")
	username := r.FormValue("username")
	
	fmt.Println(string(charJSON))

	var char models.Character
	decoder := json.NewDecoder(strings.NewReader(charJSON))
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&char)
	if err != nil {
		er := http.StatusBadRequest
		http.Error(w, err.Error(), er)
		return
	}

	err = database.InsertCharacter(char, username)
	if err != nil {
		er := http.StatusBadRequest
		http.Error(w, err.Error(), er)
		return
	}
	fmt.Fprintf(w, "character %v was succesfully added", char.Name)
}

func GetCharacters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := http.StatusMethodNotAllowed
		http.Error(w, "method not allowed", err)
		return
	}

	username := r.FormValue("username")

	userChars, err := database.GetCharactersFromUser(username)
	if err != nil {
		er := http.StatusBadRequest
		http.Error(w, err.Error(), er)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userChars)
}