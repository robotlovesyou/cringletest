// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
	"github.com/robotlovesyou/cringletest"
	"github.com/robotlovesyou/cringletest/clclient"
	"github.com/spf13/cobra"
)

// bestCmd represents the best command
var bestCmd = &cobra.Command{
	Use:   "best [from currency] to [to currency] [--address someone@example.com]",
	Short: "Get the best exchange rate between the given currencies over the last 7 days",
	Long: `
cconv best fetches the best exchange rate between two currencies from the last 7 days.
if an email address is supplied the result will be emailed in additon to being reported on the command line.

For example:

cconv best GBP to EUR --address someone@example.com

would get the best exchange rate between GBP and EUR and would mail the result to someone@example.com

cconv best ignores the --date flag`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("not enough args to rate")
		}

		if args[1] != "to" {
			return errors.New("incorrect argument format")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		from := strings.ToUpper(args[0])
		to := args[2]

		client, err := clclient.New()
		if err != nil {
			errorResult(err)
			return
		}

		notifiers, err := getNotifiers()
		if err != nil {
			errorResult(err)
			return
		}

		err = fetchBest(context.Background(), &requestConfig{
			From:      from,
			To:        []string{to},
			Client:    client,
			Notifiers: notifiers,
		})
		if err != nil {
			errorResult(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(bestCmd)
}

type rateResult struct {
	rates cringletest.RateMap
	err   error
}

// pass requestConfig by value to create a copy which can be modified safely
func getRateOnDate(ctx context.Context, config requestConfig, date time.Time) chan rateResult {
	config.Date = date
	ch := make(chan rateResult)
	go func() {
		rm, err := getRates(ctx, &config)
		ch <- rateResult{rates: rm, err: err}
	}()
	return ch
}

func getLastSevenRates(ctx context.Context, config *requestConfig) (rates []*cringletest.ExchangeRate, err error) {
	date := time.Now()
	channels := []chan rateResult{}

	for i := 7; i > 0; i-- {
		channels = append(channels, getRateOnDate(ctx, *config, date))
		date = date.Add(-24 * time.Hour)
	}

	for _, ch := range channels {
		result := <-ch
		if result.err != nil {
			return nil, errors.Wrap(result.err, "error getting best rates")
		}

		// There will be only one rate in each map so we can add it to the rates
		for _, rate := range result.rates {
			rates = append(rates, rate)
		}
	}
	return rates, nil
}

func selectBestRate(rates []*cringletest.ExchangeRate) *cringletest.ExchangeRate {
	bestRate := &cringletest.ExchangeRate{Value: new(decimal.Big)}

	for _, rate := range rates {
		if bestRate.Value.Cmp(rate.Value) < 0 {
			bestRate = rate
		}
	}

	return bestRate
}

func fetchBest(ctx context.Context, config *requestConfig) error {
	rates, err := getLastSevenRates(ctx, config)
	if err != nil {
		return err
	}

	if len(rates) == 0 {
		return cringletest.ErrBadCurrencies
	}

	bestRate := selectBestRate(rates)

	nfunc := func(cx context.Context, n cringletest.Notifier, value *decimal.Big, rates []*cringletest.ExchangeRate) error {
		return n.NotifyBest(ctx, rates[0])
	}

	return notifyAll(ctx, config.Notifiers, decimal.New(1, 0), []*cringletest.ExchangeRate{bestRate}, nfunc)
}
