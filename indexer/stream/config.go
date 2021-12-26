package stream

import "time"

type StubConfig struct {
	// all fields have defaults test purposes
	TickDuration    time.Duration `yaml:"tick_duration" default:"5s"`
	PriceMedian     float64       `yaml:"price_median" default:"100"`
	PriceDerivation float64       `yaml:"price_derivation" default:"20"`
}
