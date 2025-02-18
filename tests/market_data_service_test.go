package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"testing"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/database"
	"github.com/chyngyz-sydykov/marketpulse/internal/marketdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Test Suite Struct
type MarketDataServiceTestSuite struct {
	suite.Suite
	db      *sql.DB
	service *marketdata.MarketDataService
}

func (suite *MarketDataServiceTestSuite) SetupSuite() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	mockRedis := &MockRedisService{}

	suite.service = marketdata.NewMarketDataService(mockRedis)
	suite.db = database.DB
}

// Clear DB Before Each Test
func (suite *MarketDataServiceTestSuite) SetupTest() {
	_, err := suite.db.Exec("DELETE FROM data_btc_1h; DELETE FROM data_btc_4h; DELETE FROM data_btc_1d;")
	if err != nil {
		log.Fatalf("❌ Failed to clean DB: %v", err)
	}
}

func calculateTrend(open, close float64) float64 {
	return math.Round((close-open)/open*10000) / 10000
}

// ✅ Test Case 1: Store 1H Record
func (suite *MarketDataServiceTestSuite) TestShouldStoreSingleRecord() {
	record := &binance.RecordDto{
		Symbol:     "BTCUSDT",
		Timeframe:  "1h",
		Timestamp:  time.Now().Truncate(time.Hour),
		Open:       110,
		High:       120,
		Low:        100,
		Close:      105,
		Volume:     111,
		IsComplete: true,
	}

	err := suite.service.StoreData("btc", record)
	assert.NoError(suite.T(), err)

	// Verify stored record
	stored := suite.getDataByTimeframeAndTimestamp("1h", record.Timestamp)
	suite.assertRecordValues(*record, stored)
	assert.Equal(suite.T(), calculateTrend(record.Open, record.Close), stored.Trend)
}

// ✅ Test Case 2: Store Multiple 1H Records
func (suite *MarketDataServiceTestSuite) TestShouldStoreMultipleRecordsWithBatchUpsert() {
	records := []*binance.RecordDto{
		{Symbol: "BTCUSDT", Timeframe: "1h", Timestamp: time.Now().Add(-time.Hour), Open: 210, High: 220, Low: 200, Close: 205, Volume: 222, IsComplete: true},
		{Symbol: "BTCUSDT", Timeframe: "1h", Timestamp: time.Now(), Open: 211, High: 221, Low: 201, Close: 206, Volume: 223, IsComplete: true},
	}

	err := suite.service.UpsertBatchData("btc", records)
	assert.NoError(suite.T(), err)

	for _, rec := range records {
		stored := suite.getDataByTimeframeAndTimestamp("1h", rec.Timestamp)
		suite.assertRecordValues(*rec, stored)
		assert.Equal(suite.T(), calculateTrend(rec.Open, rec.Close), stored.Trend)
	}
}

// ✅ Test Case 3: Group Single 1H Record to 4H
func (suite *MarketDataServiceTestSuite) TestShoulGroupTwo1HRecordsToNotComplete4HourRecord() {
	timestamp1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 0, 0, 0, time.Now().Location())
	timestamp2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 2, 0, 0, 0, time.Now().Location())
	suite.db.Exec(`INSERT INTO data_btc_1h (symbol, timeframe, timestamp, open, high, low, close, volume, is_complete)
							 VALUES
							 ('btc', '1h', $1, 1, 10, 1, 5, 1, true),
							 ('btc', '1h', $2, 2, 11, 0, 6, 1, true)`, timestamp1, timestamp2)

	err := suite.service.StoreGroupedRecords("btc", "4h")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM data_btc_4h`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	var stored binance.RecordDto
	err = suite.db.QueryRow(`SELECT open, high, low, close, volume, is_complete, trend FROM data_btc_4h ORDER BY timestamp asc LIMIT 1`).
		Scan(&stored.Open, &stored.High, &stored.Low, &stored.Close, &stored.Volume, &stored.IsComplete, &stored.Trend)
	assert.NoError(suite.T(), err)
	suite.assertRecordValues(binance.RecordDto{
		Open:       1,
		High:       11,
		Low:        0,
		Close:      6,
		Volume:     2,
		IsComplete: false,
	}, stored)
	assert.Equal(suite.T(), calculateTrend(stored.Open, stored.Close), stored.Trend)
}

func (suite *MarketDataServiceTestSuite) TestShouldGroupUpdateExistingInCompleteGroupRecordWhenNew1HRecordIsAdded() {
	timestamp1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 0, 0, 0, time.Now().Location())
	timestamp2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 2, 0, 0, 0, time.Now().Location())
	suite.db.Exec(`INSERT INTO data_btc_1h (symbol, timeframe, timestamp, open, high, low, close, volume, is_complete)
							 VALUES
							 ('btc', '1h', $1, 1, 10, 1, 5, 1, true),
							 ('btc', '1h', $2, 2, 11, 0, 6, 1, true)`, timestamp1, timestamp2)

	err := suite.service.StoreGroupedRecords("btc", "4h")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM data_btc_4h`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	var stored binance.RecordDto
	err = suite.db.QueryRow(`SELECT open, high, low, close, volume, is_complete, trend FROM data_btc_4h ORDER BY timestamp asc LIMIT 1`).
		Scan(&stored.Open, &stored.High, &stored.Low, &stored.Close, &stored.Volume, &stored.IsComplete, &stored.Trend)
	assert.NoError(suite.T(), err)
	suite.assertRecordValues(binance.RecordDto{
		Open:       1,
		High:       11,
		Low:        0,
		Close:      6,
		Volume:     2,
		IsComplete: false,
	}, stored)
	assert.Equal(suite.T(), calculateTrend(stored.Open, stored.Close), stored.Trend)

	timestamp3 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 3, 0, 0, 0, time.Now().Location())
	suite.db.Exec(`INSERT INTO data_btc_1h (symbol, timeframe, timestamp, open,close, low, high, volume, is_complete)
							 VALUES ('btc', '1h', $1, 2, 9, 0, 12, 2, true)`, timestamp3)
	// update existing 4h record
	err = suite.service.StoreGroupedRecords("btc", "4h")
	assert.NoError(suite.T(), err)

	suite.db.QueryRow(`SELECT COUNT(*) FROM data_btc_4h`).Scan(&count)
	assert.Equal(suite.T(), 1, count)

	err = suite.db.QueryRow(`SELECT open, high, low, close, volume, is_complete, trend FROM data_btc_4h ORDER BY timestamp asc LIMIT 1`).
		Scan(&stored.Open, &stored.High, &stored.Low, &stored.Close, &stored.Volume, &stored.IsComplete, &stored.Trend)
	assert.NoError(suite.T(), err)
	suite.assertRecordValues(binance.RecordDto{
		Open:       1,
		Close:      9,
		Low:        0,
		High:       12,
		Volume:     4,
		IsComplete: false,
	}, stored)
	assert.Equal(suite.T(), calculateTrend(stored.Open, stored.Close), stored.Trend)
}

func (suite *MarketDataServiceTestSuite) TestShoulGroup1HRecordsWithinTimeTangeToComplete4HourRecord() {
	timestamp1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 0, 0, 0, time.Now().Location())
	timestamp2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 2, 0, 0, 0, time.Now().Location())
	timestamp3 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 3, 0, 0, 0, time.Now().Location())
	timestamp4 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 0, 0, 0, time.Now().Location())
	suite.db.Exec(`INSERT INTO data_btc_1h (symbol, timeframe, timestamp, open, close, high, low, volume, is_complete)
							 VALUES
							 ('btc', '1h', $1, 1, 5, 10, 3, 1, true),
							 ('btc', '1h', $2, 5, 6, 11, 1, 1, true),
							 ('btc', '1h', $3, 6, 7, 9, 0, 1, true),
							 ('btc', '1h', $4, 7, 8, 8, 3, 1, true)`, timestamp1, timestamp2, timestamp3, timestamp4)

	err := suite.service.StoreGroupedRecords("btc", "4h")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM data_btc_4h`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	var complete4hData binance.RecordDto
	err = suite.db.QueryRow(`SELECT open, high, low, close, volume, is_complete, trend FROM data_btc_4h ORDER BY timestamp asc LIMIT 1`).
		Scan(&complete4hData.Open, &complete4hData.High, &complete4hData.Low, &complete4hData.Close, &complete4hData.Volume, &complete4hData.IsComplete, &complete4hData.Trend)
	assert.NoError(suite.T(), err)
	suite.assertRecordValues(binance.RecordDto{
		Open:       1,
		Close:      8,
		Low:        0,
		High:       11,
		Volume:     4,
		IsComplete: true,
	}, complete4hData)
	assert.Equal(suite.T(), calculateTrend(complete4hData.Open, complete4hData.Close), complete4hData.Trend)
}

func (suite *MarketDataServiceTestSuite) TestShoulGroup1HRecordsExceedingTimeRangeToCompleteAndInComplete4HourRecords() {

	timestamp1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 0, 0, 0, time.Now().Location())
	timestamp2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 2, 0, 0, 0, time.Now().Location())
	timestamp3 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 3, 0, 0, 0, time.Now().Location())
	timestamp4 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 0, 0, 0, time.Now().Location())
	timestamp5 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 5, 0, 0, 0, time.Now().Location())
	timestamp6 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 6, 0, 0, 0, time.Now().Location())
	suite.db.Exec(`INSERT INTO data_btc_1h (symbol, timeframe, timestamp, open, close, high, low, volume, is_complete)
							 VALUES
							 ('btc', '1h', $1, 1, 5, 10, 3, 1, true),
							 ('btc', '1h', $2, 5, 6, 11, 1, 1, true),
							 ('btc', '1h', $3, 6, 7, 9, 0, 1, true),
							 ('btc', '1h', $4, 7, 8, 8, 3, 1, true),
							 ('btc', '1h', $5, 1, 6, 9, 2, 1, true),
							 ('btc', '1h', $6, 6, 1, 11, 3, 1, true)`,
		timestamp1, timestamp2, timestamp3, timestamp4, timestamp5, timestamp6)

	err := suite.service.StoreGroupedRecords("btc", "4h")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM data_btc_4h`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)

	complete4hData := suite.getDataByTimeframeAndTimestamp("4h", timestamp4)
	suite.assertRecordValues(binance.RecordDto{
		Open:       1,
		Close:      8,
		Low:        0,
		High:       11,
		Volume:     4,
		IsComplete: true,
	}, complete4hData)
	assert.Equal(suite.T(), calculateTrend(complete4hData.Open, complete4hData.Close), complete4hData.Trend)

	next4HGroupTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Now().Location())
	incomplete4hData := suite.getDataByTimeframeAndTimestamp("4h", next4HGroupTime)
	suite.assertRecordValues(binance.RecordDto{
		Open:       1,
		Close:      1,
		Low:        2,
		High:       11,
		Volume:     2,
		IsComplete: false,
	}, incomplete4hData)
	assert.Equal(suite.T(), calculateTrend(incomplete4hData.Open, incomplete4hData.Close), incomplete4hData.Trend)
}

func (suite *MarketDataServiceTestSuite) TestShoulGroup1HRecordsBetweenDaysTo4HourRecords() {
	timestamp1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 22, 0, 0, 0, time.Now().Location())
	timestamp2 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 23, 0, 0, 0, time.Now().Location())
	midnight := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 24, 0, 0, 0, time.Now().Location())
	timestamp4 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 0, 0, 0, time.Now().Location())
	timestamp5 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 2, 0, 0, 0, time.Now().Location())
	timestamp6 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 3, 0, 0, 0, time.Now().Location())
	suite.db.Exec(`INSERT INTO data_btc_1h (symbol, timeframe, timestamp, open, close, high, low, volume, is_complete)
							 VALUES
							 ('btc', '1h', $1, 100, 105, 110, 80, 10, true),
							 ('btc', '1h', $2, 105, 110, 160, 100, 15, true),
							 ('btc', '1h', $3, 110, 115, 120, 105, 20, true),
							 ('btc', '1h', $4, 115, 120, 125, 110, 25, true),
							 ('btc', '1h', $5, 120, 125, 150, 70, 30, true),
							 ('btc', '1h', $6, 125, 130, 135, 120, 35, true)`,
		timestamp1, timestamp2, midnight, timestamp4, timestamp5, timestamp6)

	err := suite.service.StoreGroupedRecords("btc", "4h")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM data_btc_4h`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)

	complete4hData := suite.getDataByTimeframeAndTimestamp("4h", midnight)

	suite.assertRecordValues(binance.RecordDto{
		Open:       100,
		Close:      115,
		Low:        80,
		High:       160,
		Volume:     45,
		IsComplete: true,
	}, complete4hData)

	next4HGroupTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 0, 0, 0, time.Now().Location())
	incomplete4hData := suite.getDataByTimeframeAndTimestamp("4h", next4HGroupTime)
	suite.assertRecordValues(binance.RecordDto{
		Open:       115,
		Close:      130,
		Low:        70,
		High:       150,
		Volume:     90,
		IsComplete: false,
	}, incomplete4hData)
	assert.Equal(suite.T(), calculateTrend(incomplete4hData.Open, incomplete4hData.Close), incomplete4hData.Trend)
}

//TODO create grouped record, then updated group record. check that data is updated

func (suite *MarketDataServiceTestSuite) TestShouldGroup1HRecordsTo1DRecords() {
	_, err := suite.db.Exec("DELETE FROM data_btc_1h; DELETE FROM data_btc_1d;")
	if err != nil {
		log.Fatalf("❌ Failed to clean DB: %v", err)
	}
	var records []*binance.RecordDto

	for i := 1; i < 50; i++ {
		tmpTime := time.Date(2020, 1, 1, i+1, 0, 0, 0, time.Now().Location())
		records = append(records, &binance.RecordDto{
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
	}

	err = suite.service.UpsertBatchData("btc", records)
	assert.NoError(suite.T(), err)

	err = suite.service.StoreGroupedRecords("btc", "1d")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM data_btc_1d where is_complete = true`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)
	err = suite.db.QueryRow(`SELECT COUNT(*) FROM data_btc_1d where is_complete = false`).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	jan2Time := time.Date(2020, 1, 1, 24, 0, 0, 0, time.Now().Location())
	jan2Data := suite.getDataByTimeframeAndTimestamp("1d", jan2Time)
	suite.assertRecordValues(binance.RecordDto{
		Open:       1,
		Close:      23,
		Low:        0,
		High:       24,
		Volume:     46,
		IsComplete: true,
	}, jan2Data)

	jan3Time := time.Date(2020, 1, 2, 24, 0, 0, 0, time.Now().Location())
	jan3Data := suite.getDataByTimeframeAndTimestamp("1d", jan3Time)
	suite.assertRecordValues(binance.RecordDto{
		Open:       24,
		Close:      47,
		Low:        23,
		High:       48,
		Volume:     48,
		IsComplete: true,
	}, jan3Data)

	jan4Time := time.Date(2020, 1, 3, 24, 0, 0, 0, time.Now().Location())
	jan4Data := suite.getDataByTimeframeAndTimestamp("1d", jan4Time)
	suite.assertRecordValues(binance.RecordDto{
		Open:       48,
		Close:      49,
		Low:        47,
		High:       50,
		Volume:     4,
		IsComplete: false,
	}, jan4Data)
}

func (suite *MarketDataServiceTestSuite) assertRecordValues(expected, actual binance.RecordDto) {
	assert.Equal(suite.T(), expected.Open, actual.Open)
	assert.Equal(suite.T(), expected.Close, actual.Close)
	assert.Equal(suite.T(), expected.Low, actual.Low)
	assert.Equal(suite.T(), expected.High, actual.High)
	assert.Equal(suite.T(), expected.Volume, actual.Volume)
	assert.Equal(suite.T(), expected.IsComplete, actual.IsComplete)
}

func (suite *MarketDataServiceTestSuite) getDataByTimeframeAndTimestamp(timeframe string, time time.Time) binance.RecordDto {
	var data binance.RecordDto
	query := fmt.Sprintf(`SELECT open, high, low, close, volume, is_complete, trend FROM data_btc_%s WHERE timestamp = $1`, timeframe)
	err := suite.db.QueryRow(query, time).
		Scan(&data.Open, &data.High, &data.Low, &data.Close, &data.Volume, &data.IsComplete, &data.Trend)
	assert.NoError(suite.T(), err)
	return data
}

// ✅ Run the Test Suite
func TestMarketDataService(t *testing.T) {
	suite.Run(t, new(MarketDataServiceTestSuite))
}
