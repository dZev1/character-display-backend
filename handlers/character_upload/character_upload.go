package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/dZev1/character-display/database"
	loginHandlers "github.com/dZev1/character-display/handlers/login"
	"github.com/dZev1/character-display/utils"
)

func UploadCharacter(w http.ResponseWriter, r *http.Request) {
	var err error

	err = loginHandlers.Authorize(r)
	if err != nil {
		http.Error(w,"unauthorized", http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	charJSON := r.FormValue("char_json")
		
	char, err := utils.JsonToChar(charJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	char.Name = cases.Title(language.Und , cases.NoLower).String(char.Name)
	

	err = database.InsertCharacter(char, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "character %v was succesfully added", char.Name)
}

func GetCharacters(w http.ResponseWriter, r *http.Request) {
	field := r.FormValue("field")
	value := r.FormValue("value")

	userChars, err := database.GetCharactersByField(field, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userChars)
}

func GetAllCharacters(w http.ResponseWriter, r *http.Request) {
	chars, err := database.GetAllCharacters()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chars)
}

func EditCharacter(w http.ResponseWriter, r *http.Request) {
	var err error
	
	username := r.FormValue("username")
	charName := r.FormValue("char_name")

	charName = cases.Title(language.Und , cases.NoLower).String(charName)

	err = loginHandlers.Authorize(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		char, err := database.GetCharacter(username, charName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(char)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	if r.Method == http.MethodPut {
		charJSON := r.FormValue("char_json")
		
		char, err := utils.JsonToChar(charJSON)
		char.Name = cases.Title(language.Und , cases.NoLower).String(char.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = database.UpdateCharacter(username, char)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "character updated successfully")
	}
}

func DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	var err error
	
	err = loginHandlers.Authorize(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	charName := r.FormValue("char_name")

	charName = cases.Title(language.Und , cases.NoLower).String(charName)
	
	err = database.DeleteCharacter(username, charName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "character %v deleted successfully", charName)
}

