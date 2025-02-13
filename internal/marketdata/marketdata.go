package marketdata

import (
	"fmt"
	"log"
	"slices"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	aggregator "github.com/chyngyz-sydykov/marketpulse/internal/marketdata/aggregator"
)

type MarketDataService struct {
	repository MarketDataRepository
	aggregator aggregator.Aggregator
}

func NewMarketDataService() *MarketDataService {
	repository := NewMarketDataRepository()
	aggregator := aggregator.NewAggregator()
	return &MarketDataService{
		repository: *repository,
		aggregator: *aggregator,
	}
}

func (service *MarketDataService) Store4HourRecords(currency string) error {

	last4HourRecord, err := service.repository.getLastRecord(currency, config.FOUR_HOUR)
	if err != nil {
		return err
	}
	//fmt.Println("last 4h record:", currency, last4HourRecord)
	var oneHourRecords []binance.RecordDto

	if last4HourRecord == nil {
		// No previous 4-hour record → Fetch all 1-hour records
		oneHourRecords, err = service.repository.getRecords(currency, config.ONE_HOUR) // Get all
	} else {
		// Fetch only 1-hour records after the last 4-hour timestamp
		oneHourRecords, err = service.repository.getRecordsAfter(currency, config.ONE_HOUR, last4HourRecord.Timestamp)
	}
	if err != nil {
		return err
	}

	if len(oneHourRecords) < 4 {
		log.Printf("Not enough 1-hour records to create a 4-hour record for %s", currency)
		return nil
	}

	var aggregatedRecords []binance.RecordDto
	for i := 0; i+3 < len(oneHourRecords); i += 4 {
		group := oneHourRecords[i : i+4]
		aggregatedRecord := service.aggregator.Aggregate(group, config.FOUR_HOUR)
		aggregatedRecords = append(aggregatedRecords, aggregatedRecord)
	}

	// Insert aggregated records into the 4-hour table
	err = service.repository.storeDataBatch(currency, aggregatedRecords)
	if err != nil {
		return err
	}
	return nil
}

func (service *MarketDataService) StoreData(currency string, data *binance.RecordDto) error {
	if !slices.Contains(config.DefaultCurrencies, currency) {
		return fmt.Errorf("unknown currency: %s", currency)
	} else {
		if !slices.Contains(config.DefaultTimeframes, data.Timeframe) {
			return fmt.Errorf("unknown time frame: %s", data.Timeframe)
		}
		exists, err := service.repository.checkIfRecordExists(currency, data.Timeframe, data.Timestamp)
		if err != nil {
			return err
		}
		if exists {
			fmt.Printf("data already exists for currency: %s, timeframe: %s, timestamp: %s\n", currency, data.Timeframe, data.Timestamp)
			return nil
		}
		return service.repository.storeData(currency, data)
	}

}
