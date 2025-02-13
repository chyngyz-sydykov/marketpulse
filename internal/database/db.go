package database

import (
	"database/sql"
	"fmt"

	"github.com/chyngyz-sydykov/marketpulse/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() error {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	return nil
}
