package characterupload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dZev1/character-display/database"
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

	username := r.FormValue("username")
	charJSON := r.FormValue("char_json")
	
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