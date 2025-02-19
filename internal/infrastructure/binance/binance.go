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
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
)

// FetchMarketData retrieves data from Binance API
func FetchTicker(currency string) (*dto.DataDto, error) {

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

	var binanceTickerData dto.BinanceTickerResponse

	fmt.Println("body:", string(body))
	if err := json.Unmarshal(body, &binanceTickerData); err != nil {
		return nil, err
	}

	return &dto.DataDto{
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
func FetchKline(currency string, interval string, limit int) ([]*dto.DataDto, error) {

	cfg := config.LoadConfig()
	client := GetHTTPClient()
	req, err := http.NewRequest("GET", cfg.BinanceBaseAPIUrl+"klines?interval="+interval+"&limit="+strconv.Itoa(limit)+"&symbol="+strings.ToUpper(currency)+"USDT", nil)
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

	records := make([]*dto.DataDto, 0, len(binanceKlineData))

	for i := 0; i < len(binanceKlineData); i++ {
		element := binanceKlineData[i]
		records = append(records, &dto.DataDto{
			Symbol:     strings.ToUpper(currency) + "USDT",
			Timestamp:  time.Unix(int64(element[0].(float64))/1000, 0).Truncate(time.Hour), // Kline open time with precision to hour
			Timeframe:  config.ONE_HOUR,
			Open:       parseFloat(element[1].(string)), // Open price
			High:       parseFloat(element[2].(string)), // High price
			Low:        parseFloat(element[3].(string)), // Low price
			Close:      parseFloat(element[4].(string)), // Close price
			Volume:     parseFloat(element[5].(string)), // Volume
			IsComplete: true,
		})
	}
	return records, nil
}
func parseFloat(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return parsed
}
