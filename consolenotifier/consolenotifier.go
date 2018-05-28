// Package consolenotifier implements cringletest.Notify for stdout
package consolenotifier

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ericlagergren/decimal"
	"github.com/robotlovesyou/cringletest"
)

const dateFormat = "Mon 02 Jan 2006"

type notifier struct {
	// set the io.Writer as a member of the notifier so that it can be
	// modified during testing
	out io.Writer
}

const (
	ratesTitle = "Exchange Rate Results on %s:"
	valueTitle = "Currency Conversion Results on %s:"
	bestTitle  = "Best Exchange Rate in the last 7 days is:"
)

// New Returns a cringletest.Notifier which sends notifications to the console
func New() cringletest.Notifier {
	return &notifier{os.Stdout}
}

func (n *notifier) writeRateLine(value *decimal.Big, rate *cringletest.ExchangeRate) error {
	_, err := fmt.Fprintf(n.out,
		"%16.4f %6s Buys %16.4f %6s\n",
		value,
		rate.From,
		new(decimal.Big).Mul(value, rate.Value),
		rate.To,
	)
	return err
}

func (n *notifier) notifyList(title string, value *decimal.Big, rates []*cringletest.ExchangeRate) error {
	if len(rates) == 0 {
		return cringletest.ErrNoRates
	}

	fmt.Fprintln(n.out, fmt.Sprintf(title, rates[0].Date.Format(dateFormat)))
	for _, rate := range rates {
		n.writeRateLine(value, rate)
	}

	return nil
}

func (n *notifier) NotifyRates(ctx context.Context, rates []*cringletest.ExchangeRate) error {
	return n.notifyList(ratesTitle, decimal.New(1, 0), rates)
}

func (n *notifier) NotifyValue(ctx context.Context, value *decimal.Big, rates []*cringletest.ExchangeRate) error {
	return n.notifyList(valueTitle, value, rates)
}

func (n *notifier) NotifyBest(ctx context.Context, rate *cringletest.ExchangeRate) error {
	fmt.Fprintln(n.out, bestTitle)
	fmt.Fprintf(n.out,
		"1.0000 %s to %.4f %s on %s\n",
		rate.From,
		rate.Value,
		rate.To,
		rate.Date.Format(dateFormat),
	)
	return nil
}
