package marketdata

import "github.com/chyngyz-sydykov/marketpulse/internal/binance"

type Aggregator struct{}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Aggregate(group []binance.RecordDto, timeFrame string) *binance.RecordDto {
	return &binance.RecordDto{
		Timeframe: timeFrame,
		Symbol:    group[0].Symbol,
		Timestamp: group[len(group)-1].Timestamp,
		High:      a.MaxHigh(group),
		Low:       a.MinLow(group),
		Close:     group[len(group)-1].Close,
		Open:      group[0].Open,
		Volume:    a.TotalVolume(group),
	}
}

func (a *Aggregator) MaxHigh(group []binance.RecordDto) float64 {
	max := group[0].High
	for _, data := range group {
		if data.High > max {
			max = data.High
		}
	}
	return max
}

func (a *Aggregator) MinLow(group []binance.RecordDto) float64 {
	min := group[0].Low
	for _, data := range group {
		if data.Low < min {
			min = data.Low
		}
	}
	return min
}

func (a *Aggregator) TotalVolume(group []binance.RecordDto) float64 {
	total := 0.0
	for _, data := range group {
		total += data.Volume
	}
	return total
}
