package binance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
)

type BinanceMockTransport struct {
	RealTransport http.RoundTripper
}

// RoundTrip intercepts the request and returns a mock response if enabled
func (m *BinanceMockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	mockBinance := config.LoadConfig().MockBinance
	if mockBinance {
		fmt.Println("🔹 Mocking Binance API Response:", req.URL.Path)
		var mockBody []byte
		if strings.Contains(req.URL.Path, "ticker") {
			// Fake response data (customize as needed)
			mockBody, _ = json.Marshal(map[string]interface{}{
				"symbol":             "FAKE_" + req.URL.Query().Get("symbol"),
				"priceChange":        "1263.77000000",
				"priceChangePercent": "1.303",
				"weightedAvgPrice":   "97577.19891976",
				"prevClosePrice":     "96985.99000000",
				"lastPrice":          "98249.77000000",
				"lastQty":            "0.00808000",
				"bidPrice":           "98249.76000000",
				"bidQty":             "1.55120000",
				"askPrice":           "98249.77000000",
				"askQty":             "4.20582000",
				"openPrice":          "97712.00000000",
				"highPrice":          "97912.00000000",
				"lowPrice":           "96866.39000000",
				"volume":             "17212.93500000",
				"quoteVolume":        "1679589982.48793910",
				"openTime":           time.Now().Add(-1*time.Hour).Unix() * 1000,
				"closeTime":          time.Now().Unix() * 1000,
				"firstId":            4541889495,
				"lastId":             4544959049,
				"count":              3069555,
			})
		} else if strings.Contains(req.URL.Path, "klines") {
			open := 100 + math.Round(rand.Float64()*100)
			close := open * (0.5 + rand.Float64()*2)
			if rand.Float64() < 0.9 {
				close = open * (0.9 + rand.Float64()*0.2) // 10% higher or lower
			} else if rand.Float64() < 0.5 {
				close = open * (0.5 + rand.Float64())
			}
			high := math.Max(open, close) * (1.1 + rand.Float64()*0.1)
			low := math.Min(open, close) * (0.9 + rand.Float64()*0.1)
			mockBody, _ = json.Marshal([][]interface{}{
				{
					time.Now().Add(-1*time.Hour).Unix() * 1000,              // Kline open time
					fmt.Sprintf("%.5f", open),                               // Open price
					fmt.Sprintf("%.5f", high),                               // High price
					fmt.Sprintf("%.5f", low),                                // Low price
					fmt.Sprintf("%.5f", close),                              // Close price
					fmt.Sprintf("%.5f", 300+math.Round(rand.Float64()*100)), // Volume
					time.Now().Unix() * 1000,                                // Kline Close time
					"31371060.16996650",                                     // Quote asset volume
					fmt.Sprintf("%.5f", math.Round(rand.Float64()*1000)),    // Number of trades
					"141.91034000",                                          // Taker buy base asset volume
					"13947601.85652800",                                     // Taker buy quote asset volume
					"0",                                                     // Unused field, ignore.
				},
			})
		} else {
			mockBody, _ = json.Marshal(map[string]interface{}{
				"symbol":      "FAKE_BTCUSDT",
				"priceChange": "1263.77000000",
			})
		}

		// Encode response as JSON
		mockResp := io.NopCloser(bytes.NewReader(mockBody))

		// Return fake HTTP response
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       mockResp,
			Header:     make(http.Header),
		}, nil
	}

	// If mocking is disabled, make a real request
	return m.RealTransport.RoundTrip(req)
}

// GetHTTPClient returns a custom client with mock support
func GetHTTPClient() *http.Client {
	return &http.Client{
		Transport: &BinanceMockTransport{RealTransport: http.DefaultTransport},
	}
}
