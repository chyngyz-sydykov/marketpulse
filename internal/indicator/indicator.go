package indicator

import (
	"log"
	"math"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
)

type Indicator struct {
}

func NewIndicator() *Indicator {
	return &Indicator{}
}

func (indicator *Indicator) CalculateAllIndicators(history []binance.RecordDto) (*binance.IndicatorDto, error) {
	// Extract close prices from historical data
	var closePrices []float64
	for _, record := range history {
		closePrices = append(closePrices, record.Close)
	}

	if len(closePrices) == 0 {
		log.Printf("%s", config.COLOR_YELLOW+"not enough data to calculate indicators"+config.COLOR_RESET)
		return nil, nil
	}

	// Compute each indicator using historical close prices
	sma := SMA(closePrices, 14)
	ema := EMA(closePrices, 14)
	stdDev := StandardDeviation(closePrices)
	lowerBollinger := sma - 2*stdDev
	upperBollinger := sma + 2*stdDev
	rsi := RSI(closePrices, 14)
	volatility := (history[len(history)-1].High - history[len(history)-1].Low) / history[len(history)-1].Close
	macd := MACD(closePrices, 12, 26)
	macdSignal := MACDSignal(closePrices, 9)

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
	ema := prices[0]
	for _, price := range prices {
		ema = (price-ema)*multiplier + ema
	}
	return ema
}

// Standard Deviation
func StandardDeviation(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	mean := SMA(prices, len(prices))
	sum := 0.0
	for _, price := range prices {
		sum += math.Pow(price-mean, 2)
	}
	return math.Sqrt(sum / float64(len(prices)))
}

// Relative Strength Index (RSI)
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

// MACD (Moving Average Convergence Divergence)
func MACD(prices []float64, shortPeriod, longPeriod int) float64 {
	return EMA(prices, shortPeriod) - EMA(prices, longPeriod)
}

// MACD Signal Line (9-period EMA of MACD)
func MACDSignal(prices []float64, signalPeriod int) float64 {
	macdValues := make([]float64, len(prices))
	for i := range prices {
		macdValues[i] = MACD(prices[:i+1], 12, 26)
	}
	return EMA(macdValues, signalPeriod)
}

func (indicator *Indicator) Trend(data *binance.RecordDto) float64 {
	return math.Round((data.Close-data.Open)/data.Open*10000) / 10000
}
