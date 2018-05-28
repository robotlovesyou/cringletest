package cringletest

import (
	"context"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
)

var (
	//ErrNoRecipient should be returned when no recipient is configured
	ErrNoRecipient = errors.New("no notify recipient")
	// ErrSendFailed should be returned when the notification failed to send
	ErrSendFailed = errors.New("notification failed to send")
	// ErrNoRates should be returned when the notifier is called with no rates
	ErrNoRates = errors.New("no rates provided")
)

//Notifier describes a service which can send a notification about the results of a query
type Notifier interface {
	NotifyRates(ctx context.Context, rates []*ExchangeRate) error
	NotifyValue(ctx context.Context, value *decimal.Big, rates []*ExchangeRate) error
	NotifyBest(ctx context.Context, rate *ExchangeRate) error
}
