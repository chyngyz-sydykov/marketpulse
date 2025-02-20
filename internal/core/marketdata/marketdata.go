package marketdata

import (
	"context"
	"log"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	aggregator "github.com/chyngyz-sydykov/marketpulse/internal/core/marketdata/aggregator"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/validator"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/redis"
)

type MarketDataService struct {
	repository MarketDataRepository
	aggregator aggregator.Aggregator
	indicator  indicator.IndicatorService
	redis      redis.RedisServiceInterface
	validator  validator.Validator
}

func NewMarketDataService(redis redis.RedisServiceInterface) *MarketDataService {
	repository := NewMarketDataRepository()
	indicator := indicator.NewIndicatorService(redis)
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

func (service *MarketDataService) StoreData(currency string, data *dto.DataDto) error {
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

	return service.publishEvent(config.EVENT_NEW_DATA_ADDED)
}

func (service *MarketDataService) UpsertBatchData(currency string, data []*dto.DataDto) error {
	timeFrame := data[0].Timeframe
	if err := service.validator.ValidateCurrencyAndTimeframe(currency, timeFrame); err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	for _, record := range data {
		record.Trend = service.indicator.Trend(record)
	}

	if err := service.repository.upsertBatchByTimeFrame(currency, timeFrame, data); err != nil {
		return err
	}

	return service.publishEvent(config.EVENT_NEW_DATA_ADDED)
}

func (service *MarketDataService) StoreGroupedRecords(currency string, groupingTimeframe string) error {
	log.Printf(config.COLOR_BLUE+"grouping records currency:%s timeframe:%s"+config.COLOR_RESET, currency, groupingTimeframe)
	if err := service.validator.ValidateCurrencyAndTimeframe(currency, groupingTimeframe); err != nil {
		return err
	}

	lastCompleteGroupRecord, err := service.repository.getLastCompleteRecord(currency, groupingTimeframe)
	if err != nil {
		return err
	}
	hoursInGroup := config.HoursByTimeframe[groupingTimeframe]

	var oneHourRecords []dto.DataDto
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
	return service.publishEvent(config.EVENT_NEW_GROUP_DATA_ADDED)
}

func (service *MarketDataService) publishEvent(eventName string) error {
	ctx := context.Background()
	err := service.redis.PublishEvent(ctx, eventName, config.SERVICE_NAME)
	if err != nil {
		return err
	}
	return nil
}
