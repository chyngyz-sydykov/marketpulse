package marketdata

import (
	"fmt"
	"log"
	"slices"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/indicator"
	aggregator "github.com/chyngyz-sydykov/marketpulse/internal/marketdata/aggregator"
)

type MarketDataService struct {
	repository MarketDataRepository
	aggregator aggregator.Aggregator
	indicator  indicator.Indicator
}

func NewMarketDataService() *MarketDataService {
	repository := NewMarketDataRepository()
	aggregator := aggregator.NewAggregator()
	indicator := indicator.NewIndicator()
	return &MarketDataService{
		repository: *repository,
		aggregator: *aggregator,
		indicator:  *indicator,
	}
}

func (service *MarketDataService) Store4HourRecords(currency string) error {

	last4HourRecord, err := service.repository.getLastRecord(currency, config.FOUR_HOUR)
	if err != nil {
		return err
	}
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
		log.Printf(config.COLOR_YELLOW+"not enough 1-hour records for %s to create a 4-hour record"+config.COLOR_RESET, currency)
		return nil
	}

	var aggregatedRecords []*binance.RecordDto
	for i := 0; i+3 < len(oneHourRecords); i += 4 {
		group := oneHourRecords[i : i+4]
		aggregatedRecord := service.aggregator.Aggregate(group, config.FOUR_HOUR)
		aggregatedRecord.Trend = service.indicator.Trend(aggregatedRecord)
		aggregatedRecords = append(aggregatedRecords, aggregatedRecord)

		// calculate indicators for a group of 4 records
		// save it to indicator table
		_, _ = service.indicator.CalculateAllIndicators(group)
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
			return fmt.Errorf("unknown timeframe: %s", data.Timeframe)
		}

		exists, err := service.repository.checkIfRecordExists(currency, data.Timeframe, data.Timestamp)
		if err != nil {
			return err
		}
		if exists {
			log.Printf(config.COLOR_YELLOW+"data already exists for %s %s %s\n"+config.COLOR_RESET, currency, data.Timeframe, data.Timestamp)
			return nil
		}

		data.Trend = service.indicator.Trend(data)

		return service.repository.storeData(currency, data)
	}

}
