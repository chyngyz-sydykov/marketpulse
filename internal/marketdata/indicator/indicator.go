package marketdata

import (
	"math"

	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
)

type Indicator struct {
}

func NewIndicator() *Indicator {
	return &Indicator{}
}

func (indicator *Indicator) CalculateAllIndicators(data *binance.RecordDto) (*binance.IndicatorDto, error) {
	//return &binance.IndicatorDto{}, nil
	prices := []float64{data.Open, data.High, data.Low, data.Close} // Example: Use actual historical prices

	// Compute each indicator sequentially
	sma := indicator.SMA(prices, 14)
	ema := indicator.EMA(prices, 14)
	stdDev := indicator.StandardDeviation(prices)
	lowerBollinger := sma - 2*stdDev
	upperBollinger := sma + 2*stdDev
	rsi := indicator.RSI(prices, 14)
	volatility := (data.High - data.Low) / data.Close
	macd := indicator.MACD(prices, 12, 26)
	macdSignal := indicator.MACDSignal(prices, 9)

	// Store results in IndicatorDto
	indicatorDto := &binance.IndicatorDto{
		SMA:            sma,
		EMA:            ema,
		StdDev:         stdDev,
		LowerBollinger: lowerBollinger,
		UpperBollinger: upperBollinger,
		RSI:            rsi,
		Volatility:     volatility,
		MACD:           macd,
		MACDSignal:     macdSignal,
	}

	return indicatorDto, nil
}

// Simple Moving Average (SMA)
func (indicator *Indicator) SMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}
	sum := 0.0
	for _, price := range prices[len(prices)-period:] {
		sum += price
	}
	return sum / float64(period)
}

// Exponential Moving Average (EMA)
func (indicator *Indicator) EMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}
	multiplier := 2.0 / (float64(period) + 1.0)
	ema := prices[0] // Initial EMA value
	for _, price := range prices {
		ema = (price-ema)*multiplier + ema
	}
	return ema
}

// Standard Deviation
func (indicator *Indicator) StandardDeviation(prices []float64) float64 {
	mean := indicator.SMA(prices, len(prices))
	var sum float64
	for _, price := range prices {
		sum += math.Pow(price-mean, 2)
	}
	return math.Sqrt(sum / float64(len(prices)))
}

// Bollinger Bands
func (indicator *Indicator) BollingerBands(prices []float64, period int) (float64, float64) {
	sma := indicator.SMA(prices, period)
	stddev := indicator.StandardDeviation(prices)
	return sma - 2*stddev, sma + 2*stddev
}

// RSI (Relative Strength Index)
func (indicator *Indicator) RSI(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}
	gain, loss := 0.0, 0.0
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
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}
func (indicator *Indicator) MACD(prices []float64, shortPeriod, longPeriod int) float64 {
	return indicator.EMA(prices, shortPeriod) - indicator.EMA(prices, longPeriod)
}

// MACD Signal Line (9-period EMA of MACD)
func (indicator *Indicator) MACDSignal(prices []float64, signalPeriod int) float64 {
	macdValues := make([]float64, len(prices))
	for i := range prices {
		macdValues[i] = indicator.MACD(prices[:i+1], 12, 26)
	}
	return indicator.EMA(macdValues, signalPeriod)
}
