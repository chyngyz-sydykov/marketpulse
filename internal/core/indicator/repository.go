package indicator

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/database"
)

type IndicatorRepository struct {
}

func NewIndicatorRepository() *IndicatorRepository {
	return &IndicatorRepository{}
}

// 1. Have no linked indicator record
// 2. OR Have is_complete = false in data table
func (repo *IndicatorRepository) GetUnprocessedMarketData(currency, timeframe string) ([]dto.DataDto, error) {
	db := database.DB

	query := fmt.Sprintf(`
		SELECT d.*
		FROM data_%s_%s d
		WHERE d.is_complete = false
		OR NOT EXISTS (
			SELECT 1 FROM indicator_%s_%s i
			WHERE i.data_timestamp = d.timestamp
		)
		ORDER BY d.timestamp ASC;
	`, currency, timeframe, currency, timeframe)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []dto.DataDto
	for rows.Next() {
		var record dto.DataDto
		err := rows.Scan(&record.Id, &record.Symbol, &record.Timeframe, &record.Timestamp, &record.Open, &record.Close, &record.High, &record.Low, &record.Volume, &record.Trend, &record.IsComplete)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (repository *IndicatorRepository) StreamOneHourRecords(ctx context.Context, currency, timeframe string) (<-chan dto.DataDto, error) {
	db := database.DB
	query := fmt.Sprintf(`
		SELECT id, symbol, timeframe, timestamp, open, close, low, high, volume, trend, is_complete FROM data_%s_%s
		ORDER BY timestamp ASC
	`, currency, timeframe)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Create a channel to stream records
	dataChan := make(chan dto.DataDto)

	go func() {
		defer close(dataChan)
		defer rows.Close()

		for rows.Next() {
			var record dto.DataDto
			err := rows.Scan(
				&record.Id,
				&record.Symbol,
				&record.Timeframe,
				&record.Timestamp,
				&record.Open,
				&record.Close,
				&record.Low,
				&record.High,
				&record.Volume,
				&record.Trend,
				&record.IsComplete,
			)
			if err == nil {
				dataChan <- record
			}
		}
	}()

	return dataChan, nil
}
func (repository *IndicatorRepository) getLastRecord(currency, timeframe string) (*dto.IndicatorDto, error) {
	query := fmt.Sprintf(`SELECT timeframe, timestamp, sma, ema, std_dev, lower_bollinger, upper_bollinger, volatility, rsi, macd, macd_signal, data_timestamp 
						  FROM indicator_%s_%s ORDER BY timestamp DESC LIMIT 1`, currency, timeframe)

	row := database.DB.QueryRow(query)

	var record dto.IndicatorDto
	err := row.Scan(&record.Timeframe, &record.Timestamp, &record.SMA, &record.EMA, &record.StdDev, &record.LowerBollinger, &record.UpperBollinger, &record.Volatility, &record.RSI, &record.MACD, &record.MACDSignal, &record.DataTimestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("💾 error fetching data: %v", err)
	}

	return &record, nil
}

func (repository *IndicatorRepository) storeData(currency string, data *dto.IndicatorDto) error {
	//TODO should this be here?
	if !slices.Contains(config.DefaultCurrencies, currency) {
		return fmt.Errorf("unknown currency: %s", currency)
	} else {
		tx, err := database.DB.Begin()
		if err != nil {
			log.Println("💾 Error starting transaction:", err)
			return err
		}

		query := `INSERT INTO indicator_` + currency + ` (timeframe, timestamp, sma, ema, std_dev, lower_bollinger, upper_bollinger, volatility, rsi, macd, macd_signal, data_timestamp) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = tx.Exec(query, data.Timeframe, data.Timestamp, data.SMA, data.EMA, data.StdDev, data.LowerBollinger, data.UpperBollinger, data.Volatility, data.RSI, data.MACD, data.MACDSignal, data.DataTimestamp)

		if err != nil {
			tx.Rollback() // Rollback transaction if insert fails
			log.Println("💾 Error inserting indicator:", err)
			return err
		}
		err = tx.Commit() // Commit transaction
		if err != nil {
			log.Println("💾 Error committing transaction:", err)
			return err
		}

		log.Println("💾 ✅ indicator stored successfully!")
		return nil
	}
}

func (repository *IndicatorRepository) upsertBatchByTimeFrame(currency string, timeFrame string, records []*dto.IndicatorDto) error {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("💾 Error starting transaction:", err)
		return err
	}
	var values []interface{}
	var placeholders []string

	for i, record := range records {
		fmt.Println("record", record.Timestamp)
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*12+1, i*12+2, i*12+3, i*12+4, i*12+5, i*12+6, i*12+7, i*12+8, i*12+9, i*12+10, i*12+11, i*12+12))
		values = append(values, record.Timeframe, record.Timestamp, record.SMA, record.EMA, record.StdDev, record.LowerBollinger, record.UpperBollinger, record.Volatility, record.RSI, record.MACD, record.MACDSignal, record.DataTimestamp)
	}

	query := fmt.Sprintf(`
		INSERT INTO indicator_%s_%s (timeframe, timestamp, sma, ema, std_dev, lower_bollinger, upper_bollinger, volatility, rsi, macd, macd_signal, data_timestamp)
		VALUES %s
		ON CONFLICT (timestamp) DO UPDATE SET
			sma = EXCLUDED.sma,
			ema = EXCLUDED.ema,
			std_dev = EXCLUDED.std_dev,
			lower_bollinger = EXCLUDED.lower_bollinger,
			upper_bollinger = EXCLUDED.upper_bollinger,
			volatility = EXCLUDED.volatility,
			rsi = EXCLUDED.rsi,
			macd = EXCLUDED.macd,
			macd_signal = EXCLUDED.macd_signal,
			data_timestamp = EXCLUDED.data_timestamp
	`, currency, timeFrame, strings.Join(placeholders, ","))

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
