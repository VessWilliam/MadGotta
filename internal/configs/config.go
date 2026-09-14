package configs

import (
	"gotta/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() models.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	return models.Config{
		DbPath:     getEnv("DB_PATH", "data/todo.db"),
		ServerPort: getEnv("SERVER_PORT", ":8080"),
	}
}
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
