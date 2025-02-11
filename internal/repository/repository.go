package repository

import (
	"fmt"
	"log"
	"slices"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/database"
)

// StoreMarketData saves Binance market data to PostgreSQL
func StoreMarketData(currency string, data *binance.RecordDto) error {
	if !slices.Contains(config.DefaultCurrencies, currency) {
		return fmt.Errorf("unknown currency: %s", currency)
	} else {

		tx, err := database.DB.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			return err
		}

		query := `INSERT INTO data_` + currency + ` (symbol, timestamp, timeframe, open, high, low, close, volume) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err = tx.Exec(query, data.Symbol, data.Timestamp, data.Timeframe, data.Open, data.High, data.Low, data.Close, data.Volume)
		if err != nil {
			tx.Rollback() // Rollback transaction if insert fails
			log.Println("Error inserting data:", err)
			return err
		}
		err = tx.Commit() // Commit transaction
		if err != nil {
			log.Println("Error committing transaction:", err)
			return err
		}

		log.Println("Market data stored successfully!")
		return nil
	}
	return nil
}
