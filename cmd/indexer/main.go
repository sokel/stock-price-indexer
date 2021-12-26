package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"stock-price-indexer/cmd"
	"stock-price-indexer/indexer"
)

func main() {
	//nolint:errcheck
	cmd.NewCmd(run).Execute()
}

func run(app cmd.AppContext) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := indexer.NewConfig(app.ConfigPath)
	if err != nil {
		return err
	}

	i := indexer.New(cfg)

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return cmd.WaitInterrupted(ctx)
	})

	wg.Go(func() error {
		return i.Serve(ctx)
	})

	if err := wg.Wait(); err != nil {
		return fmt.Errorf("termination: %s", err)
	}

	return nil
}
