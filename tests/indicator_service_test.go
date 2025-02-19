package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"testing"
	"time"

	indicator "github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/marketdata"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Test Suite Struct
type IndicatorServiceTestSuite struct {
	suite.Suite
	db                *sql.DB
	indicatorService  *indicator.IndicatorService
	marketDataService *marketdata.MarketDataService
}

type oclh struct {
	Open  float64
	Close float64
	Low   float64
	High  float64
}

func (suite *IndicatorServiceTestSuite) SetupSuite() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	mockRedis := &MockRedisService{}
	mockRedis.On("PublishEvent", mock.Anything, "NewRecordAdded", "MarketPulse").Return(nil)

	suite.indicatorService = indicator.NewIndicatorService()
	suite.marketDataService = marketdata.NewMarketDataService(mockRedis)
	suite.db = database.DB
}

// Clear DB Before Each Test
func (suite *IndicatorServiceTestSuite) SetupTest() {
	// _, err := suite.db.Exec("DELETE FROM data_btc_1h; DELETE FROM data_btc_4h; DELETE FROM data_btc_1d;")
	// if err != nil {
	// 	log.Fatalf("❌ Failed to clean DB: %v", err)
	// }
}

// ✅ Test Case 1: Store 1H Record
// func (suite *IndicatorServiceTestSuite) TestShouldReturnErrorIfCurrencyIsUnknown() {
// 	err := suite.indicatorService.ComputeAndStore("unknown currency", "4h")
// 	assert.Error(suite.T(), err)
// 	assert.ErrorContains(suite.T(), err, "unknown currency:")
// }
// func (suite *IndicatorServiceTestSuite) TestShouldReturnErrorIfTimeframeIsUnknown() {
// 	err := suite.indicatorService.ComputeAndStore("btc", "3h")
// 	assert.Error(suite.T(), err)
// 	assert.ErrorContains(suite.T(), err, "unknown timeframe:")
// }

// func (suite *IndicatorServiceTestSuite) TestShouldNotCreateIndicatorIfDataDoesNotExist() {
// 	_, errr := suite.db.Exec("DELETE FROM data_btc_1h; DELETE FROM data_btc_4h; DELETE FROM data_btc_1d;DELETE from indicator_btc;")
// 	if errr != nil {
// 		log.Fatalf("❌ Failed to clean DB: %v", errr)
// 	}
// 	err := suite.indicatorService.ComputeAndStore("btc", "4h")
// 	assert.NoError(suite.T(), err)

// 	var count int
// 	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), 0, count)
// }

// func (suite *IndicatorServiceTestSuite) TestShouldNotCreateIndicatorIfDataDoesNotExist() {
// 	_, errr := suite.db.Exec("DELETE FROM data_btc_1h; DELETE FROM data_btc_4h; DELETE FROM data_btc_1d;DELETE from indicator_btc;")
// 	if errr != nil {
// 		log.Fatalf("❌ Failed to clean DB: %v", errr)
// 	}
// 	err := suite.indicatorService.ComputeAndStore("btc", "4h")
// 	assert.NoError(suite.T(), err)

// 	var count int
// 	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), 0, count)
// }

// func (suite *IndicatorServiceTestSuite) TestShouldNotCreateIndicatorIfHistoryDataIsMissing() {
// 	_, errr := suite.db.Exec("DELETE FROM data_btc_1h; DELETE FROM data_btc_4h; DELETE FROM data_btc_1d;DELETE from indicator_btc;")
// 	if errr != nil {
// 		log.Fatalf("❌ Failed to clean DB: %v", errr)
// 	}

// 	err := suite.marketDataService.UpsertBatchData("btc", []*dto.DataDto{{
// 		Symbol:     "BTCUSDT",
// 		Timeframe:  "4h",
// 		Timestamp:  time.Date(2020, 1, 1, 1, 0, 0, 0, time.Now().Location()),
// 		Open:       5.0,
// 		Close:      8.0,
// 		Low:        1,
// 		High:       10,
// 		Volume:     10,
// 		Trend:      calculateTrend(5.0, 8.0),
// 		IsComplete: true,
// 	}})
// 	assert.NoError(suite.T(), err)
// 	err = suite.indicatorService.ComputeAndUpsertBatch("btc", "4h")
// 	assert.NoError(suite.T(), err)

// 	var count int
// 	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), 0, count)
// }

func (suite *IndicatorServiceTestSuite) TestAAAAA() {
	_, err := suite.db.Exec("DELETE FROM data_btc_1h; DELETE FROM data_btc_4h; DELETE FROM data_btc_1d;")
	if err != nil {
		log.Fatalf("❌ Failed to clean DB: %v", err)
	}
	var records []*dto.DataDto
	var ohlcList []oclh
	for i := 1; i < 4; i++ {
		tmpTime := time.Date(2020, 1, 1, i+1, 0, 0, 0, time.Now().Location())
		records = append(records, &dto.DataDto{
			Symbol:     "BTCUSDT",
			Timeframe:  "1h",
			Timestamp:  tmpTime,
			Open:       float64(i),
			Close:      float64(i),
			Low:        float64(i - 1),
			High:       float64(i + 1),
			Volume:     2,
			IsComplete: true,
		})
		ohlcList = append(ohlcList, oclh{
			Open:  float64(i),
			Close: float64(i),
			Low:   float64(i - 1),
			High:  float64(i + 1),
		})
	}

	suite.marketDataService.UpsertBatchData("btc", records)
	suite.marketDataService.StoreGroupedRecords("btc", "4h")
	// Act
	err = suite.indicatorService.ComputeAndUpsertBatch("btc", "4h")
	assert.NoError(suite.T(), err)
	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	inidcatorData := suite.getIndicatorByTimeframeAndTimestamp("4h", time.Date(2020, 1, 1, 4, 0, 0, 0, time.Now().Location()))
	lowerBollinger, upperBollinger := Bollinger(ohlcList)
	suite.assertIndicators(dto.IndicatorDto{
		SMA:            SMA(ohlcList),
		EMA:            EMA(ohlcList, 4),
		StdDev:         StandardDeviation(ohlcList),
		LowerBollinger: lowerBollinger,
		UpperBollinger: upperBollinger,
		RSI:            RSI(ohlcList),
	}, inidcatorData)
}

//TODO test case, where there are more 1h records than 4h records. when filtering, the 1h channel should not be closed

func (suite *IndicatorServiceTestSuite) getIndicatorByTimeframeAndTimestamp(timeframe string, time time.Time) dto.IndicatorDto {
	var indicator dto.IndicatorDto
	query := fmt.Sprintf(`SELECT sma, ema, std_dev, lower_bollinger, upper_bollinger, rsi FROM indicator_btc_%s WHERE timestamp = $1`, timeframe)
	err := suite.db.QueryRow(query, time).
		Scan(&indicator.SMA, &indicator.EMA, &indicator.StdDev, &indicator.LowerBollinger, &indicator.UpperBollinger, &indicator.RSI)
	assert.NoError(suite.T(), err)
	return indicator
}

func (suite *IndicatorServiceTestSuite) assertIndicators(expected, actual dto.IndicatorDto) {
	assert.Equal(suite.T(), expected.SMA, actual.SMA)
	assert.Equal(suite.T(), 0.0, actual.EMA)
	assert.Equal(suite.T(), expected.StdDev, actual.StdDev)
	assert.Equal(suite.T(), expected.LowerBollinger, actual.LowerBollinger)
	assert.Equal(suite.T(), expected.UpperBollinger, actual.UpperBollinger)
	assert.Equal(suite.T(), expected.RSI, actual.RSI)
}

func SMA(prices []oclh) float64 {
	period := len(prices)
	sum := 0.0
	for _, price := range prices {
		sum += price.Close
	}
	fmt.Println("t SMA: ", sum, sum/float64(period))
	return sum / float64(period)
}

// Exponential Moving Average (EMA)
func EMA(prices []oclh, period int) float64 {
	previousEma := SMA(prices[:len(prices)-1])

	multiplier := 2.0 / (float64(period) + 1.0)
	latestClosePrice := prices[len(prices)-1].Close
	return latestClosePrice*multiplier + previousEma*(1-multiplier)

	// previousIndicator, err := service.repository.getLastRecord(service.currency, config.FOUR_HOUR)
	// if err != nil && err != sql.ErrNoRows {
	// 	return
	// }
	// previousEma := indicatorDto.SMA

	// if previousIndicator != nil {
	// 	previousEma = previousIndicator.EMA
	// }

	// multiplier := 2.0 / (float64(period) + 1.0)
	// latestClosePrice := prices[len(prices)-1]
	// indicatorDto.EMA = latestClosePrice*multiplier + previousEma*(1-multiplier)
}

// Standard Deviation
func StandardDeviation(prices []oclh) float64 {
	mean := SMA(prices)
	sum := 0.0
	for _, price := range prices {
		sum += math.Pow(price.Close-mean, 2)
	}
	return math.Sqrt(sum / float64(len(prices)))
}

func Bollinger(prices []oclh) (float64, float64) {
	sma := SMA(prices)
	stdDev := StandardDeviation(prices)
	lowerBollinger := sma - 2*stdDev
	upperBollinger := sma + 2*stdDev
	return lowerBollinger, upperBollinger
}

func RSI(prices []oclh) float64 {
	gain, loss := 0.0, 0.0
	period := len(prices)
	for i := 1; i < period; i++ {
		change := prices[i].Close - prices[i-1].Close
		if change > 0 {
			gain += change
		} else {
			loss -= change
		}
	}
	avgGain := gain / float64(period)
	avgLoss := loss / float64(period)
	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

func TestIndicatorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(IndicatorServiceTestSuite))
}
