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

	char, err := jsonToChar(charJSON)
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
	field := r.FormValue("field")
	value := r.FormValue("value")

	userChars, err := database.GetCharactersByField(field, value)
	if err != nil {
		er := http.StatusBadRequest
		http.Error(w, err.Error(), er)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userChars)
}


func EditCharacter(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	charName := r.FormValue("char_name")

	loginHandlers.Authorize(r)

	if r.Method == http.MethodGet {
		char, err := database.GetCharacter(username, charName)
		if err != nil {
			er := http.StatusConflict
			http.Error(w, err.Error(), er)
			return
		}
		
		encoder := json.NewEncoder(w)
		err = encoder.Encode(char)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if r.Method == http.MethodPut {
		charJSON := r.FormValue("char_json")
		
		char, err := jsonToChar(charJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = database.UpdateCharacter(username, char)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}


}

func jsonToChar(charJSON string) (models.Character, error) {
	var char models.Character
	decoder := json.NewDecoder(strings.NewReader(charJSON))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&char)
	return char, err
}