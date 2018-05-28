package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/robotlovesyou/cringletest"
	"github.com/robotlovesyou/cringletest/consolenotifier"
	"github.com/robotlovesyou/cringletest/sgnotifier"
)

type requestConfig struct {
	From      string
	To        []string
	Value     *decimal.Big
	Date      time.Time
	Client    cringletest.RateClient
	Notifiers []cringletest.Notifier
}

func errorResult(err error) {
	fmt.Printf("Cannot get rates: %v\n", err)
}

func getNotifiers() ([]cringletest.Notifier, error) {
	notifiers := []cringletest.Notifier{consolenotifier.New()}
	if len(targetAddress) != 0 {
		mailNotifier, err := sgnotifier.New(targetAddress)
		if err != nil {
			return nil, err
		}
		notifiers = append(notifiers, mailNotifier)
	}

	return notifiers, nil
}

func getTargetDate() (date time.Time, err error) {
	if len(targetDate) != 0 {
		date, err = time.Parse("2006-02-01", targetDate)
	}
	return date, err
}

func getRates(ctx context.Context, config *requestConfig) (rateMap cringletest.RateMap, err error) {
	if config.Date.IsZero() {
		return config.Client.Get(ctx, config.From, config.To...)
	}
	return config.Client.GetOn(ctx, config.Date, config.From, config.To...)
}

func rateMapToSlice(rm cringletest.RateMap) (rates []*cringletest.ExchangeRate) {
	for _, rate := range rm {
		rates = append(rates, rate)
	}
	return rates
}
