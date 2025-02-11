package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables")
	} else {
		log.Println(".env file loaded successfully.")
	}

	log.Println("RUNSIGNUP_API_URL:", GetEnv("RUNSIGNUP_API_URL", "NOT SET"))
	log.Println("RUNSIGNUP_API_KEY:", GetEnv("RUNSIGNUP_API_KEY", "NOT SET"))
	log.Println("RUNSIGNUP_API_SECRET:", GetEnv("RUNSIGNUP_API_SECRET", "NOT SET"))
}

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}