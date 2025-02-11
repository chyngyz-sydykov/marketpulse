package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
)

// FetchMarketData retrieves data from Binance API
func FetchTicker(currency string) (*RecordDto, error) {

	cfg := config.LoadConfig()
	client := GetHTTPClient()
	req, err := http.NewRequest("GET", cfg.BinanceBaseAPIUrl+"ticker/24hr?symbol="+strings.ToUpper(currency)+"USDT", nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
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

	var binanceTickerData BinanceTickerResponse

	fmt.Println("body:", string(body))
	if err := json.Unmarshal(body, &binanceTickerData); err != nil {
		return nil, err
	}

	return &RecordDto{
		Symbol:    binanceTickerData.Symbol,
		Timestamp: time.Now(),
		Timeframe: config.ONE_HOUR,
		Open:      parseFloat(binanceTickerData.OpenPrice),
		High:      parseFloat(binanceTickerData.HighPrice),
		Low:       parseFloat(binanceTickerData.LowPrice),
		Close:     parseFloat(binanceTickerData.LastPrice),
		Volume:    parseFloat(binanceTickerData.Volume),
	}, nil
}

// FetchMarketData retrieves data from Binance API
func FetchKline(currency string) (*RecordDto, error) {

	cfg := config.LoadConfig()
	client := GetHTTPClient()
	req, err := http.NewRequest("GET", cfg.BinanceBaseAPIUrl+"klines?interval=1h&limit=1&symbol="+strings.ToUpper(currency)+"USDT", nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
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

	var binanceKlineData [][]interface{}
	if err := json.Unmarshal(body, &binanceKlineData); err != nil {
		return nil, err
	}

	if len(binanceKlineData) == 0 || len(binanceKlineData[0]) < 12 {
		return nil, fmt.Errorf("unexpected response format")
	}

	firstElement := binanceKlineData[0]

	return &RecordDto{
		Symbol:    strings.ToUpper(currency) + "USDT",
		Timestamp: time.Unix(int64(firstElement[0].(float64))/1000, 0), // Kline open time
		Timeframe: config.ONE_HOUR,
		Open:      parseFloat(firstElement[1].(string)), // Open price
		High:      parseFloat(firstElement[2].(string)), // High price
		Low:       parseFloat(firstElement[3].(string)), // Low price
		Close:     parseFloat(firstElement[4].(string)), // Close price
		Volume:    parseFloat(firstElement[5].(string)), // Volume
	}, nil
}
func parseFloat(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return parsed
}
