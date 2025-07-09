package internal

import (
	"log"

	"github.com/joho/godotenv"
)

func Init() {
	// load env from .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}
}
