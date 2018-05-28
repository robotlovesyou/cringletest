package cmd

import (
	"context"
	"strings"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
	"github.com/robotlovesyou/cringletest"
	"github.com/robotlovesyou/cringletest/clclient"
	"github.com/spf13/cobra"
)

// rateCmd represents the rate command
var rateCmd = &cobra.Command{
	Use:   "rate [from currency] to [to currency]... [--date 2006-02-01] [--address someone@example.com]",
	Short: "Get one or more exchange rate, optionally on a specific date",
	Long: `
cconv rate fetches the exchange rate between one or more currencies, optionally with a specific date.
if an email address is supplied the result will be emailed in additon to being reported on the command line.

For example:

cconv rate GBP to EUR CAD --date 2018-05-25 --address someone@example.com

would get the exchange rate between GBP and both EUR and CAD on the 25th of May 2018 and would mail the result to someone@example.com
	`,
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
		to := []string{}
		for _, cur := range args[2:] {
			to = append(to, strings.ToUpper(cur))
		}

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

		date, err := getTargetDate()
		if err != nil {
			errorResult(err)
			return
		}

		err = fetchAndShow(context.Background(), &requestConfig{
			From:      from,
			To:        to,
			Date:      date,
			Client:    client,
			Notifiers: notifiers,
		})
		if err != nil {
			errorResult(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rateCmd)
}

// FetchAndShow fetches the requested exchange rates and shows them via the configured Notifiers
func fetchAndShow(ctx context.Context, config *requestConfig) error {
	rates, err := getRates(ctx, config)
	if err != nil {
		return errors.Wrap(err, "could not get rates")
	}

	nfunc := func(cx context.Context, n cringletest.Notifier, value *decimal.Big, rates []*cringletest.ExchangeRate) error {
		return n.NotifyRates(cx, rates)
	}

	return notifyAll(ctx, config.Notifiers, decimal.New(1, 0), rateMapToSlice(rates), nfunc)
}
