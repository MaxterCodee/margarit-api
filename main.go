package main

import (
	"log"

	"api-margaritai/config"
	"api-margaritai/database"
	"api-margaritai/routes"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	r := routes.SetupRouter()

	log.Println("Server running on port 8080")
	r.Run(":8080")
}
