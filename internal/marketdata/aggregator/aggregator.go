package marketdata

import (
	"time"

	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/indicator"
)

type Aggregator struct {
	indicator *indicator.Indicator
}

func NewAggregator(indicator *indicator.Indicator) *Aggregator {
	return &Aggregator{
		indicator: indicator,
	}
}
func (a *Aggregator) GroupRecords(records []binance.RecordDto, hoursInGroup int, timeframe string) ([]*binance.RecordDto, error) {
	var aggregatedRecords []*binance.RecordDto
	// Group 1-hour records into complete grouped records
	tempGroup := make([]binance.RecordDto, 0, hoursInGroup)
	for i := 0; i < len(records); i++ {
		tempGroup = append(tempGroup, records[i])
		if records[i].Timestamp.Hour()%hoursInGroup == 0 {
			if len(tempGroup) > 0 {
				aggregatedRecord := a.createCompleteAggregatedRecord(tempGroup, timeframe)
				aggregatedRecords = append(aggregatedRecords, aggregatedRecord)
				tempGroup = make([]binance.RecordDto, 0, hoursInGroup)
			}
		}
	}
	// Group 1-hour records into incomplete grouped record
	if len(tempGroup) > 0 {
		aggregatedRecord := a.createIncompleteAggregatedRecord(tempGroup, timeframe)
		aggregatedRecords = append(aggregatedRecords, aggregatedRecord)
	}
	return aggregatedRecords, nil
}

func (a *Aggregator) createCompleteAggregatedRecord(records []binance.RecordDto, timeframe string) *binance.RecordDto {
	aggregatedRecord := a.aggregate(records, timeframe)
	aggregatedRecord.Trend = a.indicator.Trend(aggregatedRecord)
	aggregatedRecord.IsComplete = true
	return aggregatedRecord
}
func (a *Aggregator) createIncompleteAggregatedRecord(records []binance.RecordDto, timeframe string) *binance.RecordDto {
	aggregatedRecord := a.aggregate(records, timeframe)
	recordTimestamp := aggregatedRecord.Timestamp
	if timeframe == "4h" {
		// set the timestamp to nearest future 4h timestamp (4h, 8h, 12h, ...)
		aggregatedRecord.Timestamp = recordTimestamp.Truncate(4 * time.Hour).Add(4 * time.Hour)
	} else if timeframe == "1d" {
		// set the timestamp to nearest future 1d timestamp
		aggregatedRecord.Timestamp = recordTimestamp.Truncate(24 * time.Hour).Add(24 * time.Hour)
	}
	aggregatedRecord.Trend = a.indicator.Trend(aggregatedRecord)
	aggregatedRecord.IsComplete = false
	return aggregatedRecord
}

func (a *Aggregator) aggregate(group []binance.RecordDto, timeFrame string) *binance.RecordDto {
	return &binance.RecordDto{
		Timeframe: timeFrame,
		Symbol:    group[0].Symbol,
		Timestamp: group[len(group)-1].Timestamp,
		High:      a.maxHigh(group),
		Low:       a.minLow(group),
		Close:     group[len(group)-1].Close,
		Open:      group[0].Open,
		Volume:    a.totalVolume(group),
	}
}

func (a *Aggregator) maxHigh(group []binance.RecordDto) float64 {
	max := group[0].High
	for _, data := range group {
		if data.High > max {
			max = data.High
		}
	}
	return max
}

func (a *Aggregator) minLow(group []binance.RecordDto) float64 {
	min := group[0].Low
	for _, data := range group {
		if data.Low < min {
			min = data.Low
		}
	}
	return min
}

func (a *Aggregator) totalVolume(group []binance.RecordDto) float64 {
	total := 0.0
	for _, data := range group {
		total += data.Volume
	}
	return total
}
