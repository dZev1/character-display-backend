package utils

import (
	"encoding/json"
	"strings"

	"github.com/dZev1/character-display/models"
)

func JsonToChar(charJSON string) (models.Character, error) {
	var char models.Character
	decoder := json.NewDecoder(strings.NewReader(charJSON))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&char)
	return char, err
}