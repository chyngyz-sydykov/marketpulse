package indicator

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/validator"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/redis"
)

type IndicatorService struct {
	repository IndicatorRepository
	validator  validator.Validator
	redis      redis.RedisServiceInterface
	currency   string
}

func NewIndicatorService(redis redis.RedisServiceInterface) *IndicatorService {
	repository := NewIndicatorRepository()
	validator := validator.NewValidator()
	return &IndicatorService{
		repository: *repository,
		validator:  *validator,
		redis:      redis,
	}
}
func (service *IndicatorService) GetRecordsByRequest(indicatorRequestDto dto.IndicatorRequestDto) ([]dto.IndicatorDto, error) {
	if err := service.validator.ValidateCurrencyAndTimeframe(indicatorRequestDto.Currency, indicatorRequestDto.Timeframe); err != nil {
		return nil, err
	}
	return service.repository.GetRecordsByRequest(indicatorRequestDto)
}

func (service *IndicatorService) ComputeAndUpsertBatch(currency string, timeframe string) error {
	log.Printf(config.COLOR_BLUE+"computing indicator for currency:%s timeframe:%s"+config.COLOR_RESET, currency, timeframe)
	if err := service.validator.ValidateCurrencyAndTimeframe(currency, timeframe); err != nil {
		return err
	}
	service.currency = currency

	// 1. Get unprocessed market data
	groupRecords, err := service.repository.GetUnprocessedMarketData(currency, timeframe)
	if err != nil {
		return err
	}
	if len(groupRecords) == 0 {
		log.Printf("%s", config.COLOR_YELLOW+"not data to calculate indicators"+config.COLOR_RESET)
		return nil
	}

	hoursInGroup := config.HoursByTimeframe[timeframe]
	ctx := context.Background()

	oneHourRecordsChan, err := service.repository.StreamOneHourRecords(ctx, "btc", "1h")
	if err != nil {
		return fmt.Errorf("indicator->StreamOneHourRecords: %w", err)
	}

	var indicators []*dto.IndicatorDto
	for _, groupRecord := range groupRecords {
		filtered1HRecords := service.filter1HRecords(groupRecord, oneHourRecordsChan, hoursInGroup)
		if len(filtered1HRecords) == 0 {
			continue
		}
		fmt.Println("filtered1HRecords: ", groupRecord.Timestamp)
		indicatorDto := &dto.IndicatorDto{}
		service.setIndicatorMetadata(indicatorDto, groupRecord)
		err = service.CalculateAllIndicators(filtered1HRecords, indicatorDto)
		// TODO collect errors
		if err != nil {
			return fmt.Errorf("indicator->CalculateAllIndicators: %w", err)
		} else {
			indicators = append(indicators, indicatorDto)
		}
	}
	if len(indicators) == 0 {
		log.Printf("%s", config.COLOR_YELLOW+"no indicators to upsert"+config.COLOR_RESET)
		return nil
	}

	err = service.repository.upsertBatchByTimeFrame(currency, timeframe, indicators)
	if err != nil {
		return fmt.Errorf("indicator->upsertBatchByTimeFrame: %w", err)
	}
	return service.publishEvent("NewIndicatorAdded")

}

func (service *IndicatorService) filter1HRecords(groupRecord dto.DataDto, oneHourRecordsChan <-chan dto.DataDto, hoursInGroup int) []dto.DataDto {

	startTime := groupRecord.Timestamp.Add(-1 * time.Duration(hoursInGroup) * time.Hour)
	endTime := groupRecord.Timestamp

	var filteredRecords []dto.DataDto
	for record := range oneHourRecordsChan {
		if record.Timestamp.After(startTime) && record.Timestamp.Before(endTime.Add(time.Hour)) {
			filteredRecords = append(filteredRecords, record)
		}
		if record.Timestamp.Equal(endTime) {
			return filteredRecords
		}
	}
	return filteredRecords
}

func (service *IndicatorService) CalculateAllIndicators(history []dto.DataDto, indicatorDto *dto.IndicatorDto) error {
	// Extract close prices from historical data
	var closePrices []float64
	var highPrices []float64
	var lowPrices []float64
	for _, record := range history {
		closePrices = append(closePrices, record.Close)
		highPrices = append(highPrices, record.High)
		lowPrices = append(lowPrices, record.Low)
	}

	// Compute each indicator using historical close prices
	service.SMA(closePrices, indicatorDto)
	service.EMA(closePrices, indicatorDto, 4)
	service.StandardDeviation(closePrices, indicatorDto)
	service.Bollinger(indicatorDto)
	service.RSI(closePrices, indicatorDto)
	service.Volatility(highPrices, lowPrices, closePrices, indicatorDto)
	service.TR(history, indicatorDto)
	// fmt.Println("closePrices: ", closePrices)
	// fmt.Println("highPrices: ", highPrices)
	// fmt.Println("lowPrices: ", lowPrices)
	// fmt.Println("SMA: ", indicatorDto.SMA)
	// fmt.Println("EMA: ", indicatorDto.EMA)
	// fmt.Println("StandardDeviation: ", indicatorDto.StdDev)
	// fmt.Println("lowerBollinger: ", indicatorDto.LowerBollinger)
	// fmt.Println("upperBollinger: ", indicatorDto.UpperBollinger)
	// fmt.Println("RSI: ", indicatorDto.RSI)
	// fmt.Println("Volatility: ", indicatorDto.Volatility)
	// fmt.Println("DataTimestamp: ", indicatorDto.DataTimestamp)
	// macdSignal := indicator.MACDSignal(closePrices, 9)
	return nil
}

func (service *IndicatorService) setIndicatorMetadata(indicatorDto *dto.IndicatorDto, dataDto dto.DataDto) {
	indicatorDto.Timeframe = dataDto.Timeframe
	indicatorDto.Timestamp = dataDto.Timestamp
	indicatorDto.DataTimestamp = dataDto.Timestamp
}

func (service *IndicatorService) SMA(prices []float64, indicatorDto *dto.IndicatorDto) {
	period := len(prices)
	sum := 0.0
	for _, price := range prices {
		sum += price
	}
	indicatorDto.SMA = sum / float64(period)
}

// Exponential Moving Average (EMA)
func (service *IndicatorService) EMA(prices []float64, indicatorDto *dto.IndicatorDto, period int) {
	previousIndicator, err := service.repository.getLastRecord(service.currency, config.FOUR_HOUR) // TODO: make timeframe dynamic
	if err != nil && err != sql.ErrNoRows {
		return
	}
	previousEma := indicatorDto.SMA

	if previousIndicator != nil {
		previousEma = previousIndicator.EMA
	}

	multiplier := 2.0 / (float64(period) + 1.0)
	latestClosePrice := prices[len(prices)-1]
	indicatorDto.EMA = latestClosePrice*multiplier + previousEma*(1-multiplier)
}

// True Range (TR)
func (service *IndicatorService) TR(history []dto.DataDto, indicatorDto *dto.IndicatorDto) {
	previousData, err := service.repository.getPreviousMarketData(service.currency, config.FOUR_HOUR, indicatorDto.Timestamp)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	high := history[0].High
	low := history[0].Low
	previousClose := 0.0
	for i := 1; i < len(history); i++ {
		high = math.Max(history[i].High, high)
		low = math.Min(history[i].Low, low)
	}
	if (previousData == dto.DataDto{}) {
		previousClose = previousData.Close
	}

	indicatorDto.TR = math.Max(high-low, math.Max(math.Abs(high-previousClose), math.Abs(low-previousClose)))
}

// Standard Deviation
func (service *IndicatorService) StandardDeviation(prices []float64, indicatorDto *dto.IndicatorDto) {
	mean := indicatorDto.SMA
	sum := 0.0
	for _, price := range prices {
		sum += math.Pow(price-mean, 2)
	}
	indicatorDto.StdDev = math.Sqrt(sum / float64(len(prices)))
}

func (service *IndicatorService) Bollinger(indicatorDto *dto.IndicatorDto) {
	indicatorDto.LowerBollinger = indicatorDto.SMA - 2*indicatorDto.StdDev
	indicatorDto.UpperBollinger = indicatorDto.SMA + 2*indicatorDto.StdDev
}

func (service *IndicatorService) RSI(prices []float64, indicatorDto *dto.IndicatorDto) {
	gain, loss := 0.0, 0.0
	period := len(prices)
	for i := 1; i < period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gain += change
		} else {
			loss -= change
		}
	}
	avgGain := gain / float64(period)
	avgLoss := loss / float64(period)
	if avgLoss == 0 {
		indicatorDto.RSI = 100
	}
	rs := avgGain / avgLoss
	indicatorDto.RSI = 100 - (100 / (1 + rs))
}

// // MACD (Moving Average Convergence Divergence)
// func (service *Indicator) MACD(prices []float64, shortPeriod, longPeriod int) float64 {
// 	return indicator.EMA(prices, shortPeriod) - indicator.EMA(prices, longPeriod)
// }

// // MACD Signal Line (9-period EMA of MACD)
// func (service *Indicator) MACDSignal(prices []float64, signalPeriod int) float64 {
// 	macdValues := make([]float64, len(prices))
// 	for i := range prices {
// 		macdValues[i] = indicator.MACD(prices[:i+1], 12, 26)
// 	}
// 	return indicator.EMA(macdValues, signalPeriod)
// }

func (service *IndicatorService) Trend(data *dto.DataDto) float64 {
	return math.Round((data.Close-data.Open)/data.Open*10000) / 10000
}
func (service *IndicatorService) Volatility(highPrices []float64, lowPrices []float64, closePrices []float64, indicatorDto *dto.IndicatorDto) {
	indicatorDto.Volatility = (highPrices[len(highPrices)-1] - lowPrices[len(lowPrices)-1]) / closePrices[len(highPrices)-1]
}

func (service *IndicatorService) publishEvent(eventName string) error {
	ctx := context.Background()
	err := service.redis.PublishEvent(ctx, eventName, config.APPLICATION_NAME)
	if err != nil {
		return err
	}
	return nil
}
