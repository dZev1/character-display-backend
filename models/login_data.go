package models

type User struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	SessionToken   string `json:"session_token"`
	CSRFToken      string `json:"csrf_token"`
}
