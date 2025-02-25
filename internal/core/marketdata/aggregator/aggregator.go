package marketdata

import (
	"time"

	"github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
)

type Aggregator struct {
	indicator *indicator.IndicatorService
}

func NewAggregator(indicator *indicator.IndicatorService) *Aggregator {
	return &Aggregator{
		indicator: indicator,
	}
}
func (a *Aggregator) GroupRecords(records []dto.DataDto, hoursInGroup int, timeframe string) ([]*dto.DataDto, error) {
	var aggregatedRecords []*dto.DataDto
	// Group 1-hour records into complete grouped records
	tempGroup := make([]dto.DataDto, 0, hoursInGroup)
	for i := 0; i < len(records); i++ {
		tempGroup = append(tempGroup, records[i])
		if records[i].Timestamp.Hour()%hoursInGroup == 0 {
			if len(tempGroup) > 0 {
				aggregatedRecord := a.createCompleteAggregatedRecord(tempGroup, timeframe)
				aggregatedRecords = append(aggregatedRecords, aggregatedRecord)
				tempGroup = make([]dto.DataDto, 0, hoursInGroup)
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

func (a *Aggregator) createCompleteAggregatedRecord(records []dto.DataDto, timeframe string) *dto.DataDto {
	aggregatedRecord := a.aggregate(records, timeframe)
	aggregatedRecord.Trend = a.indicator.Trend(aggregatedRecord)
	aggregatedRecord.IsComplete = true
	return aggregatedRecord
}
func (a *Aggregator) createIncompleteAggregatedRecord(records []dto.DataDto, timeframe string) *dto.DataDto {
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

func (a *Aggregator) aggregate(group []dto.DataDto, timeFrame string) *dto.DataDto {
	return &dto.DataDto{
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

func (a *Aggregator) maxHigh(group []dto.DataDto) float64 {
	max := group[0].High
	for _, data := range group {
		if data.High > max {
			max = data.High
		}
	}
	return max
}

func (a *Aggregator) minLow(group []dto.DataDto) float64 {
	min := group[0].Low
	for _, data := range group {
		if data.Low < min {
			min = data.Low
		}
	}
	return min
}

func (a *Aggregator) totalVolume(group []dto.DataDto) float64 {
	total := 0.0
	for _, data := range group {
		total += data.Volume
	}
	return total
}
