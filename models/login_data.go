package models

type User struct {
	Username       string
	HashedPassword string
	SessionToken   string
	CSFRToken      string
}
