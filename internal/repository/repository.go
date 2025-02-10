package repository

import (
	"log"

	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/database"
)

// StoreMarketData saves Binance market data to PostgreSQL
func StoreMarketData(data *binance.MarketData) error {
	query := `INSERT INTO market_data (symbol, timestamp, open_price, high_price, low_price, close_price, volume) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := database.DB.Exec(query, data.Symbol, data.Timestamp, data.Open, data.High, data.Low, data.Close, data.Volume)
	if err != nil {
		log.Println("Error inserting market data:", err)
		return err
	}
	log.Println("Market data stored successfully!")
	return nil
}
