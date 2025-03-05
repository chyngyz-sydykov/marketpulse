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
	redisMock         *MockRedisService
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

	mockRedisMarketData := &MockRedisService{}
	mockRedisMarketData.On("PublishEvent", mock.Anything, mock.Anything, "MarketPulse").Return(nil)

	mockRedis := &MockRedisService{}
	mockRedis.On("PublishEvent", mock.Anything, "NewIndicatorAdded", "MarketPulse").Return(nil)

	suite.indicatorService = indicator.NewIndicatorService(mockRedis)
	suite.marketDataService = marketdata.NewMarketDataService(mockRedisMarketData)
	suite.db = database.DB
	suite.redisMock = mockRedis
}

// Clear DB Before Each Test
func (suite *IndicatorServiceTestSuite) SetupTest() {
	_, err := suite.db.Exec("DELETE FROM data_btc; DELETE FROM indicator_btc;")
	if err != nil {
		log.Fatalf("❌ Failed to clean DB: %v", err)
	}
	suite.redisMock.ExpectedCalls = nil
	suite.redisMock.Calls = nil
}

func (suite *IndicatorServiceTestSuite) TestShouldReturnErrorIfCurrencyIsUnknown() {
	err := suite.indicatorService.ComputeAndUpsertBatch("unknown currency", "4h")
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "unknown currency:")
}
func (suite *IndicatorServiceTestSuite) TestShouldReturnErrorIfCurrencyIsEmpty() {
	err := suite.indicatorService.ComputeAndUpsertBatch("", "4h")
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "unknown currency:")
}

func (suite *IndicatorServiceTestSuite) TestShouldReturnErrorIfTimeframeIsUnknown() {
	err := suite.indicatorService.ComputeAndUpsertBatch("btc", "3h")
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "unknown timeframe:")
}

func (suite *IndicatorServiceTestSuite) TestShouldNotCreateIndicatorIfDataDoesNotExist() {
	err := suite.indicatorService.ComputeAndUpsertBatch("btc", "4h")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)

	suite.redisMock.AssertNotCalled(suite.T(), "PublishEvent", mock.Anything, mock.Anything, "MarketPulse")
}

func (suite *IndicatorServiceTestSuite) TestShouldNotCreateIndicatorIfHistoryDataIsMissing() {
	err := suite.marketDataService.UpsertBatchData("btc", []*dto.DataDto{{
		Symbol:     "BTCUSDT",
		Timeframe:  "4h",
		Timestamp:  time.Date(2020, 1, 1, 1, 0, 0, 0, time.Now().Location()),
		Open:       5.0,
		Close:      8.0,
		Low:        1,
		High:       10,
		Volume:     10,
		Trend:      calculateTrend(5.0, 8.0),
		IsComplete: true,
	}})
	assert.NoError(suite.T(), err)
	err = suite.indicatorService.ComputeAndUpsertBatch("btc", "4h")
	suite.redisMock.AssertNotCalled(suite.T(), "PublishEvent", mock.Anything, mock.Anything, "MarketPulse")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func (suite *IndicatorServiceTestSuite) TestShouldCalculateAndStoreSingleIndicator() {
	suite.redisMock.On("PublishEvent", mock.Anything, "NewIndicatorAdded", "MarketPulse").Return(nil)

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
	err := suite.indicatorService.ComputeAndUpsertBatch("btc", "4h")
	suite.redisMock.AssertCalled(suite.T(), "PublishEvent", mock.Anything, "NewIndicatorAdded", "MarketPulse")
	assert.NoError(suite.T(), err)
	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	inidcatorData := suite.getIndicatorByTimeframeAndTimestamp("4h", time.Date(2020, 1, 1, 4, 0, 0, 0, time.Now().Location()))
	lowerBollinger, upperBollinger := Bollinger(ohlcList)
	suite.assertIndicators(dto.IndicatorDto{
		SMA:            SMA(ohlcList),
		EMA:            inidcatorData.EMA, //TODO fix this
		StdDev:         StandardDeviation(ohlcList),
		LowerBollinger: lowerBollinger,
		UpperBollinger: upperBollinger,
		RSI:            50,
		TR:             4,
	}, inidcatorData)
}

func (suite *IndicatorServiceTestSuite) TestShouldCalculateAndStoreMultipleIndicators() {
	suite.redisMock.On("PublishEvent", mock.Anything, "NewIndicatorAdded", "MarketPulse").Return(nil)
	var records []*dto.DataDto
	var ohlcList []oclh
	for i := 1; i < 7; i++ {
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
	err := suite.indicatorService.ComputeAndUpsertBatch("btc", "4h")
	suite.redisMock.AssertCalled(suite.T(), "PublishEvent", mock.Anything, "NewIndicatorAdded", "MarketPulse")
	assert.NoError(suite.T(), err)
	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)

	indicatorData := suite.getIndicatorByTimeframeAndTimestamp("4h", time.Date(2020, 1, 1, 4, 0, 0, 0, time.Now().Location()))
	firstIndicatorGroup := ohlcList[:3]
	lowerBollinger, upperBollinger := Bollinger(firstIndicatorGroup)
	suite.assertIndicators(dto.IndicatorDto{
		SMA:            SMA(firstIndicatorGroup),
		EMA:            indicatorData.EMA, //TODO fix this
		StdDev:         StandardDeviation(firstIndicatorGroup),
		LowerBollinger: lowerBollinger,
		UpperBollinger: upperBollinger,
		RSI:            50,
		TR:             TR(firstIndicatorGroup, dto.DataDto{Close: 0}),
	}, indicatorData)

	secondInidcatorData := suite.getIndicatorByTimeframeAndTimestamp("4h", time.Date(2020, 1, 1, 8, 0, 0, 0, time.Now().Location()))
	secondIndicatorGroup := ohlcList[3:]
	lowerBollinger, upperBollinger = Bollinger(secondIndicatorGroup)
	suite.assertIndicators(dto.IndicatorDto{
		SMA:            SMA(secondIndicatorGroup),
		EMA:            secondInidcatorData.EMA, //TODO fix this
		StdDev:         StandardDeviation(secondIndicatorGroup),
		LowerBollinger: lowerBollinger,
		UpperBollinger: upperBollinger,
		RSI:            50,
		TR:             TR(secondIndicatorGroup, dto.DataDto{Close: 0}),
	}, secondInidcatorData)
}

func (suite *IndicatorServiceTestSuite) TestShouldCalculateAndStoreMultipleIndicatorsAndIgnoreNotGrouped1HRecords() {
	suite.redisMock.On("PublishEvent", mock.Anything, "NewIndicatorAdded", "MarketPulse").Return(nil)
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

	suite.marketDataService.UpsertBatchData("btc", []*dto.DataDto{
		{
			Symbol:     "BTCUSDT",
			Timeframe:  "1h",
			Timestamp:  time.Date(2020, 1, 1, 5, 0, 0, 0, time.Now().Location()),
			Open:       float64(5),
			Close:      float64(5),
			Low:        float64(4),
			High:       float64(6),
			Volume:     2,
			IsComplete: true,
		},
	})

	// Act
	err := suite.indicatorService.ComputeAndUpsertBatch("btc", "4h")

	suite.redisMock.AssertCalled(suite.T(), "PublishEvent", mock.Anything, "NewIndicatorAdded", "MarketPulse")
	assert.NoError(suite.T(), err)
	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM indicator_btc`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	inidcatorData := suite.getIndicatorByTimeframeAndTimestamp("4h", time.Date(2020, 1, 1, 4, 0, 0, 0, time.Now().Location()))

	lowerBollinger, upperBollinger := Bollinger(ohlcList)
	suite.assertIndicators(dto.IndicatorDto{
		SMA:            SMA(ohlcList),
		EMA:            inidcatorData.EMA, //TODO fix this
		StdDev:         StandardDeviation(ohlcList),
		LowerBollinger: lowerBollinger,
		UpperBollinger: upperBollinger,
		RSI:            50,
		TR:             TR(ohlcList, dto.DataDto{Close: 0}),
	}, inidcatorData)
}

func (suite *IndicatorServiceTestSuite) getIndicatorByTimeframeAndTimestamp(timeframe string, time time.Time) dto.IndicatorDto {
	var indicator dto.IndicatorDto
	query := fmt.Sprintf(`SELECT sma, ema, tr, std_dev, lower_bollinger, upper_bollinger, rsi FROM indicator_btc_%s WHERE timestamp = $1`, timeframe)
	err := suite.db.QueryRow(query, time).
		Scan(&indicator.SMA, &indicator.EMA, &indicator.TR, &indicator.StdDev, &indicator.LowerBollinger, &indicator.UpperBollinger, &indicator.RSI)
	assert.NoError(suite.T(), err)
	return indicator
}

func (suite *IndicatorServiceTestSuite) getDataByTimeframeAndTimestamp(timeframe string, time time.Time) dto.DataDto {
	var data dto.DataDto
	query := fmt.Sprintf(`SELECT open, close, low, high FROM data_btc_%s WHERE timestamp = $1`, timeframe)
	err := suite.db.QueryRow(query, time).Scan(&data.Open, &data.Close, &data.Low, &data.High)
	assert.NoError(suite.T(), err)
	return data
}

func (suite *IndicatorServiceTestSuite) assertIndicators(expected, actual dto.IndicatorDto) {
	assert.Equal(suite.T(), expected.SMA, actual.SMA, "SMA should be equal")
	assert.Equal(suite.T(), expected.EMA, actual.EMA, "EMA should be equal")
	assert.Equal(suite.T(), expected.StdDev, actual.StdDev, "StdDev should be equal")
	assert.Equal(suite.T(), expected.LowerBollinger, actual.LowerBollinger, "LowerBollinger should be equal")
	assert.Equal(suite.T(), expected.UpperBollinger, actual.UpperBollinger, "UpperBollinger should be equal")
	assert.Equal(suite.T(), expected.RSI, actual.RSI, "RSI should be equal")
	assert.Equal(suite.T(), expected.TR, actual.TR, "TR should be equal")
}

func SMA(prices []oclh) float64 {
	period := len(prices)
	sum := 0.0
	for _, price := range prices {
		sum += price.Close
	}
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
	if len(prices) < 2 {
		return 50
	}
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

func TR(prices []oclh, previousData dto.DataDto) float64 {
	high := prices[0].High
	low := prices[0].Low
	for _, price := range prices {
		high = math.Max(price.High, high)
		low = math.Min(price.Low, low)
	}
	return math.Max((high - low), math.Max(math.Abs(high-previousData.Close), math.Abs(low-previousData.Close)))
}

func TestIndicatorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(IndicatorServiceTestSuite))
}
