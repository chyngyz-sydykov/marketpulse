package marketdata

import (
	"fmt"
	"log"
	"slices"
	"sync"

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

func (service *MarketDataService) StoreGroupedRecords(currency string, groupingTimeframe string) error {

	lastGroupRecord, err := service.repository.getLastRecord(currency, groupingTimeframe)
	if err != nil {
		return err
	}
	var oneHourRecords []binance.RecordDto
	hoursInGroup := config.HoursByTimeframe[groupingTimeframe]

	if lastGroupRecord == nil {
		// No previous 4-hour record → Fetch all 1-hour records
		oneHourRecords, err = service.repository.getRecords(currency, config.ONE_HOUR) // Get all
	} else {
		// Fetch only 1-hour records after the last 4-hour timestamp
		oneHourRecords, err = service.repository.getRecordsAfter(currency, config.ONE_HOUR, lastGroupRecord.Timestamp)
	}
	if err != nil {
		return err
	}

	if len(oneHourRecords) < hoursInGroup {
		log.Printf(config.COLOR_YELLOW+"not enough %s records for %s to create a %s record"+config.COLOR_RESET, config.ONE_HOUR, groupingTimeframe, currency)
		return nil
	}

	var aggregatedRecords []*binance.RecordDto
	for i := 0; i+hoursInGroup-1 < len(oneHourRecords); i += hoursInGroup {
		group := oneHourRecords[i : i+hoursInGroup]

		aggregatedRecord := service.createAggregatedRecord(group, groupingTimeframe)
		aggregatedRecords = append(aggregatedRecords, aggregatedRecord)

	}

	// Insert aggregated records into the 4-hour table
	err = service.repository.storeDataBatchByTimeFrame(currency, aggregatedRecords[0].Timeframe, aggregatedRecords)
	if err != nil {
		return err
	}
	return nil
}

func (service *MarketDataService) createAggregatedRecord(records []binance.RecordDto, timeframe string) *binance.RecordDto {
	aggregatedRecord := service.aggregator.Aggregate(records, timeframe)
	aggregatedRecord.Trend = service.indicator.Trend(aggregatedRecord)
	return aggregatedRecord
}

func (service *MarketDataService) StoreBatchData(currency string, data []*binance.RecordDto) error {
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
		return service.repository.storeDataBatchByTimeFrame(currency, timeFrame, data)
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

		return service.repository.storeData(currency, data)
	}

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
