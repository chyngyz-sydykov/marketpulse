package marketdata

import (
	"context"
	"log"
	"sync"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/indicator"
	aggregator "github.com/chyngyz-sydykov/marketpulse/internal/marketdata/aggregator"
	"github.com/chyngyz-sydykov/marketpulse/internal/redis"
	"github.com/chyngyz-sydykov/marketpulse/internal/validator"
)

type MarketDataService struct {
	repository MarketDataRepository
	aggregator aggregator.Aggregator
	indicator  indicator.Indicator
	redis      redis.RedisServiceInterface
	validator  validator.Validator
}

func NewMarketDataService(redis redis.RedisServiceInterface) *MarketDataService {
	repository := NewMarketDataRepository()
	indicator := indicator.NewIndicator()
	aggregator := aggregator.NewAggregator(indicator)
	validator := validator.NewValidator()
	return &MarketDataService{
		repository: *repository,
		aggregator: *aggregator,
		indicator:  *indicator,
		redis:      redis,
		validator:  *validator,
	}
}

func (service *MarketDataService) StoreData(currency string, data *binance.RecordDto) error {
	if err := service.validator.ValidateCurrencyAndTimeframe(currency, data.Timeframe); err != nil {
		return err
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

	err = service.repository.upsert(currency, data)
	if err != nil {
		return err
	}

	// TODO PUBLISH EVENT
	return service.publishEvent("NewRecordAdded")
}

func (service *MarketDataService) UpsertBatchData(currency string, data []*binance.RecordDto) error {
	timeFrame := data[0].Timeframe
	if err := service.validator.ValidateCurrencyAndTimeframe(currency, timeFrame); err != nil {
		return err
	}

	for _, record := range data {
		record.Trend = service.indicator.Trend(record)
	}
	return service.repository.upsertBatchByTimeFrame(currency, timeFrame, data)
}

func (service *MarketDataService) StoreGroupedRecords(currency string, groupingTimeframe string) error {
	log.Printf(config.COLOR_BLUE+"grouping records cur:%s timeframe:%s"+config.COLOR_RESET, currency, groupingTimeframe)
	if err := service.validator.ValidateCurrencyAndTimeframe(currency, groupingTimeframe); err != nil {
		return err
	}

	lastCompleteGroupRecord, err := service.repository.getLastCompleteRecord(currency, groupingTimeframe)
	if err != nil {
		return err
	}
	hoursInGroup := config.HoursByTimeframe[groupingTimeframe]

	var oneHourRecords []binance.RecordDto
	if lastCompleteGroupRecord == nil {
		// No previous 4-hour record → Fetch all 1-hour records within the last 4-hour period
		oneHourRecords, err = service.repository.getRecords(currency, config.ONE_HOUR) // Get all
	} else {
		// Fetch only 1-hour records after the last 4-hour timestamp
		oneHourRecords, err = service.repository.getCompleteRecordsAfter(currency, config.ONE_HOUR, lastCompleteGroupRecord.Timestamp)
	}
	if err != nil || len(oneHourRecords) == 0 {
		return nil
	}

	aggregatedRecords, err := service.aggregator.GroupRecords(oneHourRecords, hoursInGroup, groupingTimeframe)
	if err != nil {
		return err
	}

	if len(aggregatedRecords) != 0 {
		err = service.repository.upsertBatchByTimeFrame(currency, aggregatedRecords[0].Timeframe, aggregatedRecords)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *MarketDataService) publishEvent(eventName string) error {
	ctx := context.Background()
	err := service.redis.PublishEvent(ctx, eventName, config.SERVICE_NAME)
	if err != nil {
		return err
	}
	return nil

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
