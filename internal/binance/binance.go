package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
)

// MarketData represents the market data structure
type MarketData struct {
	Symbol    string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// BinanceResponse represents the API response structure
type BinanceResponse struct {
	Symbol    string `json:"symbol"`
	OpenPrice string `json:"openPrice"`
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
	LastPrice string `json:"lastPrice"`
	Volume    string `json:"volume"`
}

// FetchMarketData retrieves data from Binance API
func FetchMarketData() (*MarketData, error) {
	cfg := config.LoadConfig()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(cfg.BinanceAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var binanceData BinanceResponse
	if err := json.Unmarshal(body, &binanceData); err != nil {
		return nil, err
	}

	return &MarketData{
		Symbol:    binanceData.Symbol,
		Timestamp: time.Now(),
		Open:      parseFloat(binanceData.OpenPrice),
		High:      parseFloat(binanceData.HighPrice),
		Low:       parseFloat(binanceData.LowPrice),
		Close:     parseFloat(binanceData.LastPrice),
		Volume:    parseFloat(binanceData.Volume),
	}, nil
}

// parseFloat safely converts string to float64
func parseFloat(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return parsed
}
