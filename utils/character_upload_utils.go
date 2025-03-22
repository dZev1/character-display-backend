package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"character-display-server/models"
)

func JsonToChar(charJSON string) (models.Character, error) {
	var char models.Character
	decoder := json.NewDecoder(strings.NewReader(charJSON))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&char)
	return char, err
}

func GetUserCharacters(rows *sql.Rows) ([]models.Character, error) {
	var userChars []models.Character

	for rows.Next() {
		var char models.Character
		var statsJSON string

		err := rows.Scan(&char.Name, &char.Race, &statsJSON)
		if err != nil {
			return userChars, err
		}

		decoder := json.NewDecoder(strings.NewReader(statsJSON))
		err = decoder.Decode(&char.Stats)
		if err != nil {
			return userChars, err
		}

		fmt.Println(char.Name, char.Race, char.Stats)
		userChars = append(userChars, char)
	}

	if err := rows.Err(); err != nil {
		return userChars, err
	}

	return userChars, nil
}