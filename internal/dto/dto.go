package dto

import (
	"time"
)

type OHLCRequestDto struct {
	Currency  string
	Timeframe string
	StartTime *time.Time
	EndTime   *time.Time
	Limit     int32
	SortField string
	SortOrder string
}

type IndicatorRequestDto struct {
	Currency  string
	Timeframe string
	StartTime *time.Time
	EndTime   *time.Time
	Limit     int32
	SortField string
	SortOrder string
}

type DataDto struct {
	Id         *int
	Symbol     string
	Timeframe  string
	Timestamp  time.Time
	Open       float64
	Close      float64
	High       float64
	Low        float64
	Volume     float64
	Trend      float64
	IsComplete bool
}

type IndicatorDto struct {
	Id             *int
	Timeframe      string
	Timestamp      time.Time
	DataTimestamp  time.Time
	SMA            float64
	EMA            float64
	TR             float64
	StdDev         float64
	LowerBollinger float64
	UpperBollinger float64
	RSI            float64
	Volatility     float64
	MACD           float64
	MACDSignal     float64
}

type BinanceTickerResponse struct {
	Symbol    string `json:"symbol"`
	OpenPrice string `json:"openPrice"`
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
	LastPrice string `json:"lastPrice"`
	Volume    string `json:"volume"`
}
