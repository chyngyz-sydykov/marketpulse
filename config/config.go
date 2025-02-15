package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// var DefaultCurrencies = []string{"btc", "eth", "bnb", "sol", "trump"}
var DefaultCurrencies = []string{"btc"}
var DefaultTimeframes = []string{"1h", "4h", "1d"}
var HoursByTimeframe = map[string]int{
	"1h":  1,
	"4h":  4,
	"1d":  24,
	"7d":  168,
	"30d": 30 * 24, // TODO assuming 1 month is 30 days
}
var ONE_HOUR = "1h"
var FOUR_HOUR = "4h"
var ONE_DAY = "1d"
var SEVEN_DAY = "7d"
var THIRTY_DAY = "30d"

var COLOR_RESET = "\033[0m"
var COLOR_RED = "\033[31m"
var COLOR_GREEN = "\033[32m"
var COLOR_YELLOW = "\033[33m"
var COLOR_BLUE = "\033[34m"
var COLOR_MAGENTA = "\033[35m"
var COLOR_CYAN = "\033[36m"
var COLOR_GRAY = "\033[37m"
var COLOR_WHITE = "\033[97m"

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
