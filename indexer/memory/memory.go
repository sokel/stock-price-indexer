package memory

import (
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"

	"stock-price-indexer/entities"
)

type Memory struct {
	mem map[string]map[int64][]float64
}

func New() *Memory {
	return &Memory{
		mem: map[string]map[int64][]float64{},
	}
}

func (m *Memory) AppendNewValue(ticker entities.TickerPrice) error {
	pf, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse price")
	}

	roundedTime := m.getRoundedMinute(ticker.Time)

	if _, ok := m.mem[ticker.Ticker]; !ok {
		m.mem[ticker.Ticker] = make(map[int64][]float64)
	}

	m.mem[ticker.Ticker][roundedTime] = append(m.mem[ticker.Ticker][roundedTime], pf)

	return nil
}

// getRoundedMinute - returns Unix nano time format rounded to minute at UTC loc.
// NOTE: this function has receiver for future cases to configure loc.
func (m *Memory) getRoundedMinute(t time.Time) int64 {
	tt := t.In(time.UTC)

	rounded := time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), 0, 0, time.UTC)

	return rounded.Unix()
}

func (m *Memory) GetIndexedValueAtTime(indexFunc, ticker string, t time.Time) (int64, float64, error) {
	roundedTime := m.getRoundedMinute(t)

	// TODO: mind about empty set at start of the minute
	if len(m.mem[ticker][roundedTime]) == 0 {
		roundedTime -= roundedTime
	}

	var result float64
	var err error

	switch indexFunc {
	case "mean":
		result, err = stats.Mean(m.mem[ticker][roundedTime])
	case "median":
		result, err = stats.Median(m.mem[ticker][roundedTime])
	case "max":
		result, err = stats.Max(m.mem[ticker][roundedTime])
	case "min":
		result, err = stats.Min(m.mem[ticker][roundedTime])
	default:
		err = errors.New("index func is not recognized")
	}
	return roundedTime, result, err
}
