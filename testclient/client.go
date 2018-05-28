// Package testclient implements a fake cringletest.RateClient for use in testing
package testclient

import (
	"context"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/robotlovesyou/cringletest"
)

type client struct {
	err error
}

// New returns a new test RateClient
func New(err error) cringletest.RateClient {
	return &client{err}
}

func makeRates(date time.Time, from string, to ...string) cringletest.RateMap {
	rates := cringletest.RateMap{}
	for _, currency := range to {
		rates[currency] = &cringletest.ExchangeRate{
			From:  from,
			To:    currency,
			Date:  date,
			Value: decimal.New(1, 0),
		}
	}
	return rates
}

func (c *client) Get(ctx context.Context, from string, to ...string) (cringletest.RateMap, error) {
	if c.err != nil {
		return nil, c.err
	}

	return makeRates(time.Now(), from, to...), nil
}

func (c *client) GetOn(ctx context.Context, date time.Time, from string, to ...string) (cringletest.RateMap, error) {
	if c.err != nil {
		return nil, c.err
	}

	return makeRates(date, from, to...), nil
}
