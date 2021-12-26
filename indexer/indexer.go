package indexer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"stock-price-indexer/entities"
	"stock-price-indexer/indexer/memory"
	"stock-price-indexer/indexer/stream"
)

type PriceStreamSubscriber interface {
	SubscribePriceStream(ctx context.Context, ticker string) <-chan entities.TickerPrice
}

type tickerStream interface {
	Subscribe(ctx context.Context, ticker string) <-chan entities.TickerPrice
}

type Indexer struct {
	cfg *indexerConfig

	stream tickerStream
	mem    *memory.Memory
}

func New(cfg *Config) *Indexer {
	return &Indexer{
		cfg: &cfg.Indexer,

		// need to inject stream and mem
		stream: stream.NewStubStream(cfg.StubStream),
		mem:    memory.New(),
	}
}

func (i *Indexer) Serve(ctx context.Context) error {
	ch := i.SubscribePriceStream(ctx, i.cfg.Ticker)

	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-ch:
			if m.Err != nil {
				return errors.Wrap(m.Err, "getting error from channel")
			}
			fmt.Printf("Time - %s \t Price - %s  \r\n", m.Time, m.Price)
		}
	}

}

func (i *Indexer) SubscribePriceStream(ctx context.Context, ticker string) <-chan entities.TickerPrice {
	ch := make(chan entities.TickerPrice)
	go i.start(ctx, ticker, ch)
	return ch
}

func (i *Indexer) start(ctx context.Context, ticker string, ch chan<- entities.TickerPrice) {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return i.streamsRoutine(ctx, ticker)
	})

	wg.Go(func() error {
		return i.indexRoutine(ctx, ch, ticker)
	})

	if err := wg.Wait(); err != nil {
		ch <- entities.TickerPrice{
			Err: err,
		}
	}
}

func (i *Indexer) indexRoutine(ctx context.Context, ch chan<- entities.TickerPrice, ticker string) error {
	log.Println(i.cfg.TickerPeriod)
	markTk := time.NewTicker(i.cfg.TickerPeriod)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-markTk.C:
			log.Println("ir tc")
			resultTime, result, err := i.mem.GetIndexedValueAtTime(i.cfg.IndexFunc, ticker, time.Now())
			if err != nil {
				return errors.Wrap(err, "failed to get index price")
			}
			ch <- entities.TickerPrice{
				Ticker: ticker,
				Time:   time.Unix(resultTime, 0),
				Price:  fmt.Sprintf("%f", result),
				Err:    nil,
			}
		}
	}
}

func (i *Indexer) streamsRoutine(ctx context.Context, ticker string) error {
	streamCh := i.stream.Subscribe(ctx, ticker)

	for {
		select {
		case <-ctx.Done():
			return nil
		case newTicker := <-streamCh:
			if newTicker.Err != nil {
				return errors.Wrap(newTicker.Err, "stream returns unexpected error")
			}

			err := i.mem.AppendNewValue(newTicker)
			if err != nil {
				return errors.Wrap(err, "failed to add new ticker to memory")
			}
		}
	}
}
