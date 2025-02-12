package marketdata

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"
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

func (repository *MarketDataRepository) getRecords(currency string, timeframe string) ([]binance.RecordDto, error) {
	query := fmt.Sprintf(`SELECT id, symbol, timeframe, timestamp, open, high, low, close, volume 
						  FROM data_%s_%s ORDER BY timestamp ASC`, currency, timeframe)

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer rows.Close()

	var records []binance.RecordDto
	for rows.Next() {
		var record binance.RecordDto
		err := rows.Scan(&record.Id, &record.Symbol, &record.Timeframe, &record.Timestamp,
			&record.Open, &record.High, &record.Low, &record.Close, &record.Volume)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return records, nil
}
func (repository *MarketDataRepository) getRecordsAfter(currency string, timeframe string, lastTime time.Time) ([]binance.RecordDto, error) {
	query := fmt.Sprintf(`SELECT id, symbol, timeframe, timestamp, open, high, low, close, volume 
						  FROM data_%s_%s WHERE timestamp > $1 ORDER BY timestamp ASC`, currency, timeframe)

	rows, err := database.DB.Query(query, lastTime)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer rows.Close()

	var records []binance.RecordDto
	for rows.Next() {
		var record binance.RecordDto
		err := rows.Scan(&record.Id, &record.Symbol, &record.Timeframe, &record.Timestamp,
			&record.Open, &record.High, &record.Low, &record.Close, &record.Volume)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return records, nil
}

func (repository *MarketDataRepository) getLastRecord(currency, timeframe string) (*binance.RecordDto, error) {
	query := fmt.Sprintf(`SELECT id, symbol, timeframe, timestamp, open, high, low, close, volume 
						  FROM data_%s_%s ORDER BY timestamp DESC LIMIT 1`, currency, timeframe)

	row := database.DB.QueryRow(query)

	var record binance.RecordDto
	err := row.Scan(&record.Id, &record.Symbol, &record.Timeframe, &record.Timestamp,
		&record.Open, &record.High, &record.Low, &record.Close, &record.Volume)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return &record, nil
}
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

func (repository *MarketDataRepository) storeDataBatch(currency string, records []binance.RecordDto) error {
	if !slices.Contains(config.DefaultCurrencies, currency) {
		return fmt.Errorf("unknown currency: %s", currency)
	} else if len(records) == 0 {
		return nil
	} else {

		tx, err := database.DB.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			return err
		}

		tableName := fmt.Sprintf("data_%s", currency)

		// Construct bulk insert SQL query
		var placeholders []string
		var values []interface{}

		for i, record := range records {
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8))
			values = append(values, record.Symbol, record.Timeframe, record.Timestamp, record.Open, record.High, record.Low, record.Close, record.Volume)
		}

		query := fmt.Sprintf(`
		INSERT INTO %s (symbol, timeframe, timestamp, open, high, low, close, volume) 
		VALUES %s 
		`, tableName, strings.Join(placeholders, ", "))
		_, err = tx.Exec(query, values...)
		if err != nil {
			tx.Rollback() // Rollback transaction if insert fails
			log.Println("Error inserting batch:", err)
			return err
		}
		err = tx.Commit() // Commit transaction
		if err != nil {
			log.Println("Error committing transaction:", err)
			return err
		}

		log.Println("data batch stored successfully!")
		return nil
	}
}
