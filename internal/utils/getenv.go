package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load("/home/agust/work/iron-stream/backend/.env")
	if err != nil {
		log.Fatal(err.Error())
	}
	return os.Getenv(key)
}
