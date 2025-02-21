package grpc

import (
	"context"
	"log"
	"slices"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/dto"
	pb "github.com/chyngyz-sydykov/marketpulse/proto/marketpulse"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IndicatorHandler struct {
	IndicatorService indicator.IndicatorService
}

func NewIndicatorHandler(IndicatorService indicator.IndicatorService) *IndicatorHandler {
	return &IndicatorHandler{
		IndicatorService: IndicatorService,
	}
}

func (handler *IndicatorHandler) GetIndicators(ctx context.Context, request *pb.IndicatorRequest) (*pb.IndicatorResponse, error) {
	requestDto := handler.mapRequestToDTO(request)
	records, err := handler.IndicatorService.GetRecordsByRequest(requestDto)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return nil, status.Error(codes.InvalidArgument, "resource value(s) is invalid")
	}

	indicatorResponse := handler.mapDtoToResponse(records)
	return indicatorResponse, nil
}
func (handler *IndicatorHandler) mapRequestToDTO(req *pb.IndicatorRequest) dto.IndicatorRequestDto {
	dto := dto.IndicatorRequestDto{
		Currency:  req.Currency,
		Timeframe: req.Timeframe,
	}

	if req.StartTime != nil {
		startTime := req.StartTime.AsTime()
		dto.StartTime = &startTime
	}

	if req.EndTime != nil {
		endTime := req.EndTime.AsTime()
		dto.EndTime = &endTime
	}

	if req.Limit > 0 {
		dto.Limit = req.Limit
	} else {
		dto.Limit = int32(config.DEFAULT_LIMIT)
	}

	if req.SortField != "" && slices.Contains([]string{"timestamp"}, req.SortField) {
		dto.SortField = req.SortField
	} else {
		dto.SortField = config.DEFAULT_DATA_REQUEST_SORT_FIELD
	}
	if req.SortOrder != "" && (req.SortOrder == "ASC" || req.SortOrder == "DESC") {
		dto.SortOrder = req.SortOrder
	} else {
		dto.SortOrder = config.DEFAULT_DATA_REQUEST_SORT_ORDER
	}

	return dto
}

func (handler *IndicatorHandler) mapDtoToResponse(records []dto.IndicatorDto) *pb.IndicatorResponse {
	indicatorRecords := make([]*pb.IndicatorData, len(records))
	for i, record := range records {
		indicatorRecords[i] = &pb.IndicatorData{
			Timeframe:      record.Timeframe,
			Timestamp:      timestamppb.New(record.Timestamp),
			Sma:            record.SMA,
			Ema:            record.EMA,
			StdDev:         record.StdDev,
			LowerBollinger: record.LowerBollinger,
			UpperBollinger: record.UpperBollinger,
			Rsi:            record.RSI,
			Volatility:     record.Volatility,
			Macd:           record.MACD,
			MacdSignal:     record.MACDSignal,
		}
	}
	return &pb.IndicatorResponse{
		Indicators: indicatorRecords,
	}
}
