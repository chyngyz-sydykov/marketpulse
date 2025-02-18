package indicator

import (
	"database/sql"
	"fmt"
	"log"
	"math"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
	"github.com/chyngyz-sydykov/marketpulse/internal/pkg/utils"
)

type Indicator struct {
	repository IndicatorRepository
	currency   string
}

func NewIndicator() *Indicator {
	repository := NewIndicatorRepository()
	return &Indicator{
		repository: *repository,
	}
}
func (service *Indicator) ComputeAndStoreByTimeframe(currency string, groupingTimeframe string) error {
	return nil
}

func (service *Indicator) ComputeAndStore(currency string, history []dto.RecordDto) error {

	service.currency = currency
	fmt.Println("💡 Computing indicators...")
	if len(history) == 0 {
		log.Printf("%s", config.COLOR_YELLOW+"not enough data to calculate indicators"+config.COLOR_RESET)
		return nil
	}
	indicatorDto := &dto.IndicatorDto{}
	err := service.setMetadata(indicatorDto, history)
	if err != nil {
		return err
	}
	exists, err := service.repository.checkIfRecordExistsByTimestampAndTimeframe(currency, indicatorDto)
	if err != nil {
		return err
	}
	if exists {
		log.Printf(config.COLOR_YELLOW+"indicator already exists for %s %s %s\n"+config.COLOR_RESET, currency, indicatorDto.Timeframe, indicatorDto.Timestamp)
		return nil
	}

	if err != nil {
		return fmt.Errorf("indicator->setMetadata: %w", err)
	}
	err = service.CalculateAllIndicators(history, indicatorDto)
	if err != nil {
		return fmt.Errorf("indicator->CalculateAllIndicators: %w", err)
	}

	err = service.repository.storeData(currency, indicatorDto)
	if err != nil {
		return fmt.Errorf("indicator->storeData: %w", err)
	}
	return nil
}

func (service *Indicator) CalculateAllIndicators(history []dto.RecordDto, indicatorDto *dto.IndicatorDto) error {
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

func (service *Indicator) setMetadata(indicatorDto *dto.IndicatorDto, history []dto.RecordDto) error {
	nextTimeframe := utils.GetNextTimeframe(history[0].Timeframe)
	if nextTimeframe == "" {
		return fmt.Errorf("no next timeframe for %s", history[0].Timeframe)
	}

	lastTimestamp := history[len(history)-1].Timestamp
	indicatorDto.Timeframe = nextTimeframe
	indicatorDto.Timestamp = lastTimestamp
	indicatorDto.DataTimestamp = lastTimestamp
	return nil
}

func (service *Indicator) SMA(prices []float64, indicatorDto *dto.IndicatorDto) {
	period := len(prices)
	sum := 0.0
	for _, price := range prices {
		sum += price
	}
	indicatorDto.SMA = sum / float64(period)
}

// Exponential Moving Average (EMA)
func (service *Indicator) EMA(prices []float64, indicatorDto *dto.IndicatorDto, period int) {
	previousIndicator, err := service.repository.getLastRecord(service.currency, config.FOUR_HOUR)
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

// Standard Deviation
func (service *Indicator) StandardDeviation(prices []float64, indicatorDto *dto.IndicatorDto) {
	mean := indicatorDto.SMA
	sum := 0.0
	for _, price := range prices {
		sum += math.Pow(price-mean, 2)
	}
	indicatorDto.StdDev = math.Sqrt(sum / float64(len(prices)))
}

func (service *Indicator) Bollinger(indicatorDto *dto.IndicatorDto) {
	indicatorDto.LowerBollinger = indicatorDto.SMA - 2*indicatorDto.StdDev
	indicatorDto.UpperBollinger = indicatorDto.SMA + 2*indicatorDto.StdDev
}

func (service *Indicator) RSI(prices []float64, indicatorDto *dto.IndicatorDto) {
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

func (service *Indicator) Trend(data *dto.RecordDto) float64 {
	return math.Round((data.Close-data.Open)/data.Open*10000) / 10000
}
func (service *Indicator) Volatility(highPrices []float64, lowPrices []float64, closePrices []float64, indicatorDto *dto.IndicatorDto) {
	indicatorDto.Volatility = (highPrices[len(highPrices)-1] - lowPrices[len(lowPrices)-1]) / closePrices[len(highPrices)-1]
}
