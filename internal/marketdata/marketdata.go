package marketdata

import (
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
)

type MarketDataService struct {
	repository MarketDataRepository
}

func NewMarketDataService() *MarketDataService {
	repository := NewMarketDataRepository()
	return &MarketDataService{
		repository: *repository,
	}
}

func (service *MarketDataService) StoreMarketData(currency string, data *binance.RecordDto) error {
	return service.repository.StoreData(currency, data)
}
