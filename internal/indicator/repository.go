package indicator

import (
	"database/sql"
	"fmt"
	"log"
	"slices"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/database"
)

type IndicatorRepository struct {
}

func NewIndicatorRepository() *IndicatorRepository {
	return &IndicatorRepository{}
}
func (repository *IndicatorRepository) checkIfRecordExistsByTimestampAndTimeframe(currency string, data *binance.IndicatorDto) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM indicator_%s_%s WHERE timestamp = $1)`, currency, data.Timeframe)

	var exists bool
	err := database.DB.QueryRow(query, data.Timestamp).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Record doesn't exist
		}
		log.Printf("💾 Error checking record existence: %v", err)
		return false, err
	}

	return exists, nil
}
func (repository *IndicatorRepository) getLastRecord(currency, timeframe string) (*binance.IndicatorDto, error) {
	query := fmt.Sprintf(`SELECT timeframe, timestamp, sma, ema, std_dev, lower_bollinger, upper_bollinger, volatility, rsi, macd, macd_signal, data_timestamp 
						  FROM indicator_%s_%s ORDER BY timestamp DESC LIMIT 1`, currency, timeframe)

	row := database.DB.QueryRow(query)

	var record binance.IndicatorDto
	err := row.Scan(&record.Timeframe, &record.Timestamp, &record.SMA, &record.EMA, &record.StdDev, &record.LowerBollinger, &record.UpperBollinger, &record.Volatility, &record.RSI, &record.MACD, &record.MACDSignal, &record.DataTimestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("💾 error fetching data: %v", err)
	}

	return &record, nil
}

func (repository *IndicatorRepository) storeData(currency string, data *binance.IndicatorDto) error {
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

func (repository *IndicatorRepository) storeDataBatch(currency string, records []*binance.RecordDto) error {
	return nil
}
