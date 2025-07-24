package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"character-display-server/models"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("could not verify connection: %v", err)
	}

	fmt.Println("connection to database has been established.")
	return nil
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func InsertCharacter(character models.Character, username string) error {
	var statsJSON bytes.Buffer
	encoder := json.NewEncoder(&statsJSON)
	err := encoder.Encode(character.Stats)
	if err != nil {
		return fmt.Errorf("could not encode json: %v", err)
	}

	query := `
		INSERT INTO	characters(username, name, race, stats, image)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = db.Exec(query, username, character.Name, character.Race, statsJSON.String(), character.Image)
	if err != nil {
		return fmt.Errorf("could not execute statement: %v", err)
	}
	return nil
}

func InsertUser(username, hashedPassword string) error {
	query := `
		INSERT INTO users(username, hashed_password)
		VALUES ($1, $2)
	`

	_, err := db.Exec(query, username, hashedPassword)
	if err != nil {
		return fmt.Errorf("could not execute statement: %v", err)
	}
	return nil
}

func DeleteCharacter(username, charName string) error {
	query := `
		DELETE FROM characters
		WHERE username = $1 AND name = $2
	`

	_, err := db.Exec(query, username, charName)
	if err != nil {
		return fmt.Errorf("could not find character: %v", err)
	}
	return nil
}

func GetUser(username string) (models.User, error) {
	user := models.User{}
	query := `
		SELECT username, hashed_password, session_token, csrf_token
		FROM users 
		WHERE username = $1;
	`
	err := db.QueryRow(query, username).Scan(&user.Username, &user.HashedPassword, &user.SessionToken, &user.CSRFToken)
	if err != nil {
		return user, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func GetAllCharacters() ([]models.Character, error) {
	query := `
		SELECT username, name, race, stats, image
		FROM characters
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not fetch data: %v", err)
	}
	defer rows.Close()

	userChars, err := GetUserCharacters(rows)

	return userChars, err
}

func GetCharactersByField(field, value string) ([]models.Character, error) {
	allowedFields := map[string]bool{
		"username": true,
		"name":     true,
		"race":     true,
	}

	if !allowedFields[field] {
		return nil, fmt.Errorf("field not allowed")
	}

	query := fmt.Sprintf(`
		SELECT username, name, race, stats, image
		FROM characters
		WHERE %s = $1
	`, field)

	rows, err := db.Query(query, value)
	if err != nil {
		return nil, fmt.Errorf("could not fetch data: %v", err)
	}
	defer rows.Close()

	userChars, err := GetUserCharacters(rows)

	return userChars, err
}

func GetCharacter(username, charName string) (models.Character, error) {
	var ret models.Character
	var statsJSON string
	query := `
		SELECT username, name, race, stats, image
		FROM characters
		WHERE username = $1 AND name = $2
	`
	err := db.QueryRow(query, username, charName).Scan(&ret.Username, &ret.Name, &ret.Race, &statsJSON, &ret.Image)
	if err != nil {
		return ret, err
	}
	decoder := json.NewDecoder(strings.NewReader(statsJSON))
	err = decoder.Decode(&ret.Stats)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func UpdateCookies(user models.User) error {
	query := `
		UPDATE users
		SET session_token = $1, csrf_token = $2
		WHERE username = $3
	`

	_, err := db.Exec(query, user.SessionToken, user.CSRFToken, user.Username)
	if err != nil {
		return fmt.Errorf("could not update cookies: %v", err)
	}
	return nil
}

func UpdateCharacter(username string, character models.Character) error {
	var statsJSON bytes.Buffer
	encoder := json.NewEncoder(&statsJSON)
	err := encoder.Encode(character.Stats)
	if err != nil {
		return fmt.Errorf("could not encode stats: %v", err)
	}

	query := `
		UPDATE characters
		SET race = $1, stats = $2, image = $3
		WHERE username = $4 AND name = $5
	`

	_, err = db.Exec(query, character.Race, statsJSON.String(), character.Image, username, character.Name)
	if err != nil {
		return fmt.Errorf("could not update character: %v", err)
	}

	return nil
}

func IsInDatabase(username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}
	return exists, nil
}

func GetUserCharacters(rows *sql.Rows) ([]models.Character, error) {
	var userChars []models.Character

	for rows.Next() {
		var char models.Character
		var statsJSON string

		err := rows.Scan(&char.Username, &char.Name, &char.Race, &statsJSON, &char.Image)
		if err != nil {
			return userChars, err
		}

		decoder := json.NewDecoder(strings.NewReader(statsJSON))
		err = decoder.Decode(&char.Stats)
		if err != nil {
			return userChars, err
		}
		userChars = append(userChars, char)
	}

	if err := rows.Err(); err != nil {
		return userChars, err
	}

	return userChars, nil
}
