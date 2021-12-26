package stream

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"stock-price-indexer/entities"
)

type StubStream struct {
	cfg StubConfig
}

func NewStubStream(cfg StubConfig) *StubStream {
	return &StubStream{
		cfg: cfg,
	}
}

func (s *StubStream) Subscribe(ctx context.Context, ticker string) <-chan entities.TickerPrice {
	ch := make(chan entities.TickerPrice)

	// could be more control.
	go s.start(ctx, ticker, ch)

	return ch
}

func (s *StubStream) start(ctx context.Context, ticker string, ch chan<- entities.TickerPrice) {
	tk := time.NewTicker(s.cfg.TickDuration)

	for {
		select {
		case <-ctx.Done():
			close(ch)
			log.Println("StubStream closed")
			return
		case <-tk.C:
			ch <- s.stubValue(ticker)
		}
	}
}

func (s *StubStream) stubValue(ticker string) entities.TickerPrice {
	value, err := s.randValue()
	return entities.TickerPrice{
		Ticker: ticker,
		Time:   time.Now(),
		Price:  value,
		Err:    err,
	}
}

const sal = 10000

func (s *StubStream) randValue() (string, error) {
	randRange := big.NewInt(int64(sal * s.cfg.PriceDerivation))

	twoSidedRand, err := rand.Int(rand.Reader, randRange)
	if err != nil {
		return "", errors.Wrap(err, "failed to get random int")
	}

	less := twoSidedRand.Sub(twoSidedRand, big.NewInt(sal))

	lessFloat := float64(less.Int64())

	add := lessFloat / sal

	return fmt.Sprintf("%f", add+s.cfg.PriceMedian), nil
}
