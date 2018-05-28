package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
	"github.com/robotlovesyou/cringletest"
	"github.com/robotlovesyou/cringletest/clclient"
	"github.com/spf13/cobra"
)

// valueCmd represents the value command
var valueCmd = &cobra.Command{
	Use:   "value 1.234 [from currency] to [to currency]... [--date 2006-02-01] [--address someone@example.com]",
	Short: "Get the value of the given amount when converted to one or more target currencies",
	Long: `
cconv value fetches the value of the given amount of one currency when converted to one or more currencies, optionally with a specific date.
if an email address is supplied the result will be emailed in additon to being reported on the command line.

For example:

cconv value 200 GBP to EUR CAD --date 2018-05-25 --address someone@example.com

would get result of converting 200 GPB to both EUR and CAD on the 25th of May 2018. It would mail the result to someone@example.com
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 4 {
			return errors.New("not enough args to value")
		}

		if args[2] != "to" {
			return errors.New("incorrect argument format")
		}

		if _, ok := new(decimal.Big).SetString(args[0]); !ok {
			return fmt.Errorf("%s cannot be formatted as a number", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		value, _ := new(decimal.Big).SetString(args[0])
		from := strings.ToUpper(args[1])
		to := []string{}
		for _, cur := range args[3:] {
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

		err = fetchAndConvert(context.Background(), &requestConfig{
			From:      from,
			To:        to,
			Date:      date,
			Client:    client,
			Notifiers: notifiers,
			Value:     value,
		})
		if err != nil {
			errorResult(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(valueCmd)
}

func fetchAndConvert(ctx context.Context, config *requestConfig) error {
	rates, err := getRates(ctx, config)
	if err != nil {
		return errors.Wrap(err, "could not get values")
	}

	nfunc := func(cx context.Context, n cringletest.Notifier, value *decimal.Big, rates []*cringletest.ExchangeRate) error {
		return n.NotifyValue(cx, value, rates)
	}

	return notifyAll(ctx, config.Notifiers, config.Value, rateMapToSlice(rates), nfunc)
}
