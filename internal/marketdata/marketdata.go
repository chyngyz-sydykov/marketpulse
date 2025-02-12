package marketdata

import (
	"fmt"
	"slices"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	marketdata "github.com/chyngyz-sydykov/marketpulse/internal/marketdata/indicator"
)

type MarketDataService struct {
	repository MarketDataRepository
	indicator  marketdata.Indicator
}

func NewMarketDataService() *MarketDataService {
	repository := NewMarketDataRepository()
	indicator := marketdata.NewIndicator()
	return &MarketDataService{
		repository: *repository,
		indicator:  *indicator,
	}
}

func (service *MarketDataService) StoreDataWithIndicator(currency string, data *binance.RecordDto) error {
	_, err := service.indicator.CalculateAllIndicators(data)
	//service.repository.StoreData(currency, data)
	return err
}

func (service *MarketDataService) StoreMarketData(currency string, data *binance.RecordDto) error {
	if !slices.Contains(config.DefaultCurrencies, currency) {
		return fmt.Errorf("unknown currency: %s", currency)
	} else {
		if !slices.Contains(config.DefaultTimeframes, data.Timeframe) {
			return fmt.Errorf("unknown time frame: %s", data.Timeframe)
		}
		exists, err := service.repository.checkIfRecordExists(currency, data.Timeframe, data.Timestamp)
		if err != nil {
			return err
		}
		if exists {
			fmt.Printf("data already exists for currency: %s, timeframe: %s, timestamp: %s\n", currency, data.Timeframe, data.Timestamp)
			return nil
		}
		return service.repository.storeData(currency, data)
	}

}
