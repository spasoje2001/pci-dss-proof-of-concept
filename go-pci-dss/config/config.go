package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnectionString string
}

func LoadConfig() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	return Config{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}
}
