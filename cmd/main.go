package main

import (
	"iron-stream/api"
	"iron-stream/internal/database"
	"log"
)

func main() {
	database.ConnectDB()
	app := api.Setup()
	log.Fatal(app.Listen(":8080"))
}
