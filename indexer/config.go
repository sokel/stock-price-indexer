package indexer

import (
	"time"

	"github.com/jinzhu/configor"

	"stock-price-indexer/indexer/stream"
)

type indexerConfig struct {
	Ticker       string        `yaml:"ticker" default:"BTC_USDT"`
	IndexFunc    string        `yaml:"index_func" default:"median"`
	TickerPeriod time.Duration `yaml:"ticker_period" default:"1m"`
}

type Config struct {
	Indexer    indexerConfig     `yaml:"Indexer"`
	StubStream stream.StubConfig `yaml:"stub_stream"`
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	err := configor.Load(cfg, path)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
