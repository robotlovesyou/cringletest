// Package testnotifier implements a fake cringletest.Notifier for testing
package testnotifier

import (
	"context"

	"github.com/ericlagergren/decimal"
	"github.com/robotlovesyou/cringletest"
)

type notifier struct {
	err error
}

// New returns a test cringletest.Notifier
func New(err error) cringletest.Notifier {
	return &notifier{err}
}

func (n *notifier) NotifyRates(ctx context.Context, rates []*cringletest.ExchangeRate) error {
	return n.err
}

func (n *notifier) NotifyValue(ctx context.Context, value *decimal.Big, rates []*cringletest.ExchangeRate) error {
	return n.err
}

func (n *notifier) NotifyBest(ctx context.Context, rate *cringletest.ExchangeRate) error {
	return n.err
}
