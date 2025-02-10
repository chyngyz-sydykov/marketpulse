package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config structure to hold application configuration
type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	BinanceAPI string
}

// LoadConfig reads from .env and loads it into Config struct
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	return &Config{
		AppPort:    getEnv("APPLICATION_PORT", "8000"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_DATABASE", "db"),
		DBUser:     getEnv("DB_USERNAME", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		BinanceAPI: getEnv("BINANCE_API_URL", "https://api.binance.com/api/v3/ticker/24hr?symbol=BTCUSDT"),
	}
}

// getEnv retrieves the environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
