package metrics

import "math"

// Calculate OHLC (dummy for now)
func CalculateOHLC(open, high, low, close []float64) (float64, float64, float64, float64) {
	return open[0], high[0], low[0], close[len(close)-1]
}

// Simple Moving Average (SMA)
func SMA(prices []float64, period int) float64 {
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
func EMA(prices []float64, period int) float64 {
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
func StandardDeviation(prices []float64) float64 {
	mean := SMA(prices, len(prices))
	var sum float64
	for _, price := range prices {
		sum += math.Pow(price-mean, 2)
	}
	return math.Sqrt(sum / float64(len(prices)))
}

// Bollinger Bands
func BollingerBands(prices []float64, period int) (float64, float64) {
	sma := SMA(prices, period)
	stddev := StandardDeviation(prices)
	return sma - 2*stddev, sma + 2*stddev
}

// RSI (Relative Strength Index)
func RSI(prices []float64, period int) float64 {
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
