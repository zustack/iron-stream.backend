package main

import (
	"iron-stream/api"
	"iron-stream/internal/database"
	"log"
)

func main() {
	database.ConnectDB("DB_DEV_PATH")
	app := api.Setup()
	log.Fatal(app.Listen(":8081"))
}
