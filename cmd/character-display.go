package main

import (
	"fmt"
	"net/http"

	"character-display-server/config"
	"character-display-server/database"
	"character-display-server/middleware"
	"character-display-server/routes"

	"github.com/charmbracelet/log"
)

func main() {
	var err error

	connStr, err := config.ReadConnStrEnv()
	if err != nil {
		log.Fatal(err)
	}

	port, err := config.ReadPortEnv()
	if err != nil {
		log.Fatal(err)
	}
	port = fmt.Sprintf(":%v", port) 
	fmt.Printf("starting server at http://localhost%v/\n", port)

	err = database.InitDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()
	
	router := routes.SetupRouter()
	
	if err := http.ListenAndServe(port, middleware.CORSMiddleware(router)); err != nil {
		log.Fatal(err);
	}
}