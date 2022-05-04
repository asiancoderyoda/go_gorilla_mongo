package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvFromKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
