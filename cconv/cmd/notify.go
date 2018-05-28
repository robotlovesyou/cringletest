package cmd

import (
	"context"

	"github.com/ericlagergren/decimal"
	"github.com/robotlovesyou/cringletest"
)

func notify(ctx context.Context, notifier cringletest.Notifier, value *decimal.Big, rates []*cringletest.ExchangeRate, nfunc notifyFunc) chan error {
	ch := make(chan error)
	go func() {
		ch <- nfunc(ctx, notifier, value, rates)
	}()
	return ch
}

func notifyAll(ctx context.Context, notifiers []cringletest.Notifier, value *decimal.Big, rates []*cringletest.ExchangeRate, nfunc notifyFunc) error {
	_, ok := ctx.Deadline()
	var cancel context.CancelFunc
	if !ok {
		ctx, cancel = context.WithTimeout(ctx, cringletest.NotifyTimeout)
	}
	defer cancel()

	channels := []chan error{}
	for _, notifier := range notifiers {
		channels = append(channels, notify(ctx, notifier, value, rates, nfunc))
	}

	for _, ch := range channels {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-ch:
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type notifyFunc func(context.Context, cringletest.Notifier, *decimal.Big, []*cringletest.ExchangeRate) error
