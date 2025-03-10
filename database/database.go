package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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

func IsInDatabase(username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}
	return exists, nil
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
		SELECT username, hashed_password, session_token, csfr_token
		FROM users 
		WHERE username = $1;
	`
	err := db.QueryRow(query, username).Scan(&user.Username, &user.HashedPassword, &user.SessionToken, &user.CSFRToken)
	if err != nil {
		return user, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func UpdateCookies(user models.User) {
	query := `
		UPDATE users
		SET session_token=$1, csfr_token=$2
		WHERE username=$3
	`

	_, err := db.Exec(query, user.SessionToken, user.CSFRToken, user.Username)
	if err != nil {
		log.Fatalf("could not update cookies: %v", err)
	}
}