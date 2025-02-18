package validator

import (
	"fmt"
	"slices"

	"github.com/chyngyz-sydykov/marketpulse/config"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateCurrencyAndTimeframe(currency string, timeframe string) error {
	if !slices.Contains(config.DefaultCurrencies, currency) {
		return fmt.Errorf("unknown currency: %s", currency)
	} else {
		if !slices.Contains(config.DefaultTimeframes, timeframe) {
			return fmt.Errorf("unknown timeframe: %s", timeframe)
		}
	}
	return nil
}
