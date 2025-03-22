package main

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"character-display-server/config"
	"character-display-server/database"
	"character-display-server/routes"
)

func main() {
	var err error
	fmt.Println("starting server at http://localhost:8080/")

	connStr, err := config.ReadConnStrEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()
	
	router := routes.SetupRouter()
	
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err);
	}
	
}