package utils

import (
	"github.com/chyngyz-sydykov/marketpulse/config"
)

func GetNextTimeframe(timeframe string) string {
	for i, tf := range config.DefaultTimeframes {
		if tf == timeframe && i < len(config.DefaultTimeframes)-1 {
			return config.DefaultTimeframes[i+1]
		}
	}
	return ""
}
