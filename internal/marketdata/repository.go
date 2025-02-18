package marketdata

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

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
		return nil, fmt.Errorf("💾 error fetching data: %v", err)
	}
	defer rows.Close()

	var records []binance.RecordDto

	for rows.Next() {
		var record binance.RecordDto
		err := rows.Scan(&record.Id, &record.Symbol, &record.Timeframe, &record.Timestamp,
			&record.Open, &record.High, &record.Low, &record.Close, &record.Volume)
		if err != nil {
			return nil, fmt.Errorf("💾 error scanning row: %v", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("💾 error iterating rows: %v", err)
	}

	return records, nil
}

func (repository *MarketDataRepository) getCompleteRecordsAfter(currency string, timeframe string, lastTime time.Time) ([]binance.RecordDto, error) {
	query := fmt.Sprintf(`SELECT id, symbol, timeframe, timestamp, open, high, low, close, volume 
						  FROM data_%s_%s 
						  WHERE timestamp > $1 and is_complete = true
						  ORDER BY timestamp ASC`, currency, timeframe)

	rows, err := database.DB.Query(query, lastTime)
	if err != nil {
		return nil, fmt.Errorf("💾 error fetching data: %v", err)
	}
	defer rows.Close()

	var records []binance.RecordDto

	for rows.Next() {
		var record binance.RecordDto
		err := rows.Scan(&record.Id, &record.Symbol, &record.Timeframe, &record.Timestamp,
			&record.Open, &record.High, &record.Low, &record.Close, &record.Volume)
		if err != nil {
			return nil, fmt.Errorf("💾 error scanning row: %v", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("💾 error iterating rows: %v", err)
	}

	return records, nil
}

func (repository *MarketDataRepository) getLastCompleteRecord(currency, timeframe string) (*binance.RecordDto, error) {
	query := fmt.Sprintf(`SELECT id, symbol, timeframe, timestamp, open, high, low, close, volume 
						  FROM data_%s_%s 
						  WHERE is_complete = true
						  ORDER BY timestamp DESC LIMIT 1`, currency, timeframe)

	row := database.DB.QueryRow(query)

	var record binance.RecordDto
	err := row.Scan(&record.Id, &record.Symbol, &record.Timeframe, &record.Timestamp,
		&record.Open, &record.High, &record.Low, &record.Close, &record.Volume)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("💾 error fetching data: %v", err)
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
		return false, fmt.Errorf("💾 error checking record existence: %v", err)
	}

	return exists, nil
}

func (repository *MarketDataRepository) upsert(currency string, data *binance.RecordDto) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("💾 error starting transaction: %v", err)
	}
	query := fmt.Sprintf(
		`INSERT INTO data_%s_%s (symbol, timestamp, timeframe, open, high, low, close, volume, trend, is_complete) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (timestamp) 
		DO UPDATE 
		SET symbol = EXCLUDED.symbol, 
		timeframe = EXCLUDED.timeframe, 
		open = EXCLUDED.open, 
		high = EXCLUDED.high, 
		low = EXCLUDED.low, 
		close = EXCLUDED.close, 
		volume = EXCLUDED.volume, 
		trend = EXCLUDED.trend, 
		is_complete = EXCLUDED.is_complete`, currency, data.Timeframe)

	_, err = tx.Exec(query, data.Symbol, data.Timestamp, data.Timeframe, data.Open, data.High, data.Low, data.Close, data.Volume, data.Trend, data.IsComplete)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("💾 error upserting data: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("💾 error committing transaction: %v", err)
	}

	log.Println("✅ data upserted successfully!")
	return nil

}

func (repository *MarketDataRepository) upsertBatchByTimeFrame(currency string, timeFrame string, records []*binance.RecordDto) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("💾 error starting transaction: %v", err)
	}

	var values []interface{}
	var placeholders []string

	for i, record := range records {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*10+1, i*10+2, i*10+3, i*10+4, i*10+5, i*10+6, i*10+7, i*10+8, i*10+9, i*10+10))
		values = append(values, record.Symbol, record.Timeframe, record.Timestamp,
			record.Open, record.High, record.Low, record.Close, record.Volume, record.Trend, record.IsComplete)
	}

	query := fmt.Sprintf(`
		INSERT INTO data_%s_%s (symbol, timeframe, timestamp, open, high, low, close, volume, trend, is_complete) 
		VALUES %s
		ON CONFLICT (timestamp) DO UPDATE 
		SET open = EXCLUDED.open,
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			close = EXCLUDED.close,
			volume = EXCLUDED.volume,
			trend = EXCLUDED.trend,
			is_complete = EXCLUDED.is_complete;`,
		currency, timeFrame, strings.Join(placeholders, ","))

	_, err = tx.Exec(query, values...)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("💾 error batch upsert: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("💾 error committing transaction: %v", err)
	}

	log.Println("💾 ✅ upsert batch successfully!")
	return nil
}
