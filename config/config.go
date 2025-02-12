package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DefaultCurrencies = []string{"btc", "eth", "bnb", "sol", "trump"}
var DefaultTimeframes = []string{"1h", "4h", "1d", "1w", "1m"}
var ONE_HOUR = "1h"
var FOUR_HOUR = "4h"
var ONE_DAY = "1d"
var ONE_WEEK = "1w"
var ONE_MONTH = "1m"

// Config structure to hold application configuration
type Config struct {
	AppPort           string
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
	MockBinance       bool
	BinanceBaseAPIUrl string
}

// LoadConfig reads from .env and loads it into Config struct
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	return &Config{
		AppPort:           getEnv("APPLICATION_PORT", "8000"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBName:            getEnv("DB_DATABASE", "db"),
		DBUser:            getEnv("DB_USERNAME", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "password"),
		MockBinance:       os.Getenv("MOCK_BINANCE") == "true",
		BinanceBaseAPIUrl: getEnv("BINANCE_BASE_API_URL", "https://api.binance.com/api/v3/"),
	}
}

// getEnv retrieves the environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
