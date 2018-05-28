package cringletest

import (
	"context"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
)

var (
	// ErrBadFromCurrency should be returned when the from currency does not have a rate
	ErrBadFromCurrency = errors.New("from currency does not exist")
	// ErrNoToCurrencies should be returned when no to currencies are supplied
	ErrNoToCurrencies = errors.New("no \"to\" currencies")
	// ErrBadCurrencies should be returned when no currency conversions could be made
	ErrBadCurrencies = errors.New("bad currencies")
)

// ExchangeRate describes An exchange rate of Value between From and To on Date
type ExchangeRate struct {
	From  string
	To    string
	Date  time.Time
	Value *decimal.Big
}

// RateMap is a map from string to exchange rate
type RateMap map[string]*ExchangeRate

// RateClient describes a client able to get exchange rates
type RateClient interface {
	Get(ctx context.Context, from string, to ...string) (rates RateMap, err error)
	GetOn(ctx context.Context, date time.Time, from string, to ...string) (rates RateMap, err error)
}
