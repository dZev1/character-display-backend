package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dZev1/character-display/models"

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
		INSERT INTO	characters(username, name, race, stats)
		VALUES ($1, $2, $3, $4)
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, character.Name, character.Race, statsJSON.String())
	if err != nil {
		return fmt.Errorf("could not execute statement: %v", err)
	}
	return nil
}

func InsertUser(username string, hashedPassword string) error {
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

func GetCharactersByField(field, value string) ([]models.Character, error) {
	allowedFields := map[string]bool{
		"username": true,
		"name" : true,
		"race":     true,
		"class":    true,
	}

	if !allowedFields[field] {
		return nil, fmt.Errorf("field not allowed")
	}

	query := fmt.Sprintf(`
		SELECT name, race, stats
		FROM characters
		WHERE %s = $1
	`, field)

	rows, err := db.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userChars []models.Character

	for rows.Next() {
		var char models.Character
		var statsJSON string

		err := rows.Scan(&char.Name, &char.Race, &statsJSON)
		if err != nil {
			return userChars, nil
		}
		decoder := json.NewDecoder(strings.NewReader(statsJSON))
		err = decoder.Decode(&char.Stats)
		if err != nil {
            return userChars, err
        }
		fmt.Println(char.Name, char.Race, char.Stats)
		userChars = append(userChars, char)
	}
	if err = rows.Err(); err != nil {
		return userChars, err
	}
	return userChars, nil
}

func GetCharacter(username, charName string) (models.Character, error) {
	var ret models.Character
	var statsJSON string
	query := `
		SELECT name, race, stats FROM characters
		WHERE username = $1 AND name = $2
	`
	err := db.QueryRow(query, username, charName).Scan(&ret.Name, &ret.Race, &statsJSON)
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
	query := `
		UPDATE characters
		SET race = $1, stats $2
		WHERE username = $3 AND name = $4
	`

	_, err := db.Exec(query, character.Race, character.Stats, username, character.Name)
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