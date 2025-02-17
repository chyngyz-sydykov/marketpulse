package marketdata

import (
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

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

		return service.repository.store(currency, data)
	}

}

func (service *MarketDataService) UpsertBatchData(currency string, data []*binance.RecordDto) error {
	if !slices.Contains(config.DefaultCurrencies, currency) {
		return fmt.Errorf("unknown currency: %s", currency)
	} else {
		timeFrame := data[0].Timeframe
		if !slices.Contains(config.DefaultTimeframes, timeFrame) {
			return fmt.Errorf("unknown timeframe: %s", timeFrame)
		}

		for _, record := range data {
			record.Trend = service.indicator.Trend(record)
		}
		return service.repository.upsertBatchByTimeFrame(currency, timeFrame, data)
	}
}

func (service *MarketDataService) StoreGroupedRecords(currency string, groupingTimeframe string) error {

	lastCompleteGroupRecord, err := service.repository.getLastCompleteRecord(currency, groupingTimeframe)
	if err != nil {
		return err
	}
	var oneHourRecords []binance.RecordDto
	hoursInGroup := config.HoursByTimeframe[groupingTimeframe]

	if lastCompleteGroupRecord == nil {
		// No previous 4-hour record → Fetch all 1-hour records within the last 4-hour period
		oneHourRecords, err = service.repository.getRecordsWithTimeModulo(currency, config.ONE_HOUR, hoursInGroup) // Get all
	} else {
		// Fetch only 1-hour records after the last 4-hour timestamp
		oneHourRecords, err = service.repository.getCompleteRecordsAfter(currency, config.ONE_HOUR, lastCompleteGroupRecord.Timestamp)
	}
	if err != nil {
		return err
	}

	if len(oneHourRecords) < hoursInGroup {
		// aggregate incomplete 4-hour record
		aggregatedRecord := service.createIncompleteAggregatedRecord(oneHourRecords, groupingTimeframe)
		service.repository.upsert(currency, aggregatedRecord)
		return nil
	}

	var aggregatedRecords []*binance.RecordDto
	// Group 1-hour records into 4-hour complete records
	for i := 0; i+hoursInGroup-1 < len(oneHourRecords); i += hoursInGroup {
		group := oneHourRecords[i : i+hoursInGroup]
		aggregatedRecord := service.createCompleteAggregatedRecord(group, groupingTimeframe)
		aggregatedRecords = append(aggregatedRecords, aggregatedRecord)
	}
	// If there are remaining 1-hour records, create an incomplete 4-hour record
	if (len(aggregatedRecords)*hoursInGroup < len(oneHourRecords)) && len(aggregatedRecords) != 0 {
		group := oneHourRecords[len(aggregatedRecords)*hoursInGroup:]
		aggregatedRecord := service.createIncompleteAggregatedRecord(group, groupingTimeframe)
		aggregatedRecords = append(aggregatedRecords, aggregatedRecord)
	}
	fmt.Println("aggregatedRecords", len(aggregatedRecords), aggregatedRecords[0])
	// Insert aggregated records into the 4-hour table
	if len(aggregatedRecords) != 0 {
		err = service.repository.upsertBatchByTimeFrame(currency, aggregatedRecords[0].Timeframe, aggregatedRecords)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *MarketDataService) createCompleteAggregatedRecord(records []binance.RecordDto, timeframe string) *binance.RecordDto {
	aggregatedRecord := service.aggregator.Aggregate(records, timeframe)
	aggregatedRecord.Trend = service.indicator.Trend(aggregatedRecord)
	aggregatedRecord.IsComplete = true
	return aggregatedRecord
}

func (service *MarketDataService) createIncompleteAggregatedRecord(records []binance.RecordDto, timeframe string) *binance.RecordDto {
	aggregatedRecord := service.aggregator.Aggregate(records, timeframe)
	recordTimestamp := aggregatedRecord.Timestamp
	if timeframe == "4h" {
		// set the timestamp to nearest future 4h timestamp (3h, 7h, 11h, ...)
		aggregatedRecord.Timestamp = recordTimestamp.Truncate(4 * time.Hour).Add(3 * time.Hour)
	} else if timeframe == "1d" {
		// set the timestamp to nearest future 1d timestamp
		aggregatedRecord.Timestamp = recordTimestamp.Truncate(24 * time.Hour).Add(23 * time.Hour)
	}
	aggregatedRecord.Trend = service.indicator.Trend(aggregatedRecord)
	aggregatedRecord.IsComplete = false
	return aggregatedRecord
}

func (service *MarketDataService) calculateAndStoreIndicators(currency string, groupedRecords [][]binance.RecordDto) error {
	var wg sync.WaitGroup
	for i := 0; i < len(groupedRecords); i++ {
		wg.Add(1)
		go func(group []binance.RecordDto) {
			defer wg.Done()

			err := service.indicator.ComputeAndStore(currency, group)
			if err != nil {
				log.Printf(config.COLOR_RED+"Error computing indicators for %s: %v"+config.COLOR_RESET, currency, err)
			}
		}(groupedRecords[i])
	}
	wg.Wait()
	return nil
}
