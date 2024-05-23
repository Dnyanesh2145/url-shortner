package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var enverr = godotenv.Load(".env")

func EnvDBURI(key string) string {
	if enverr != nil {
		log.Fatal("Error while loading .env file", enverr.Error())
	}
	return os.Getenv(key)
}
