package marketdata

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/database"
)

type MarketDataRepository struct {
}

func NewMarketDataRepository() *MarketDataRepository {
	return &MarketDataRepository{}
}

// StoreMarketData saves Binance market data to PostgreSQL
func (repository *MarketDataRepository) checkIfRecordExists(currency, timeframe string, dateTime time.Time) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM data_%s_%s WHERE timestamp=$1)`, currency, timeframe)

	var exists bool
	err := database.DB.QueryRow(query, dateTime).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Record doesn't exist
		}
		log.Printf("Error checking record existence: %v", err)
		return false, err
	}

	return exists, nil
}

func (repository *MarketDataRepository) getDataByTimeFrameAndTimestamp(currency string, timeframe string, dateTime time.Time) ([]binance.RecordDto, error) {
	query := `SELECT id, symbol, timeframe, timestamp, open, high, low, close, volume 
			  FROM data_` + currency + `_` + timeframe + ` WHERE timestamp=$1`

	row := database.DB.QueryRow(query, dateTime)

	var record binance.RecordDto
	err := row.Scan(&record.Id, &record.Symbol, &record.Timeframe, &record.Timestamp,
		&record.Open, &record.High, &record.Low, &record.Close, &record.Volume)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no data found for %s %s at %s", currency, timeframe, dateTime)
		}
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return []binance.RecordDto{record}, nil
}
func (repository *MarketDataRepository) storeData(currency string, data *binance.RecordDto) error {
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
}
