package consolenotifier

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/robotlovesyou/cringletest"
	"github.com/stretchr/testify/require"
)

func getTestNotifier() (cringletest.Notifier, *bytes.Buffer) {
	buf := bytes.NewBuffer(nil)
	sender := New()
	sender.(*notifier).out = buf
	return sender, buf
}

func getTestOutput(buf *bytes.Buffer) (string, error) {
	outBytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}

	return string(outBytes), nil
}

func getFormattedDate(date time.Time) string {
	return date.Format(dateFormat)
}

func TestNotifyRatesSendsOK(t *testing.T) {
	r := require.New(t)
	sender, buf := getTestNotifier()

	rates := []*cringletest.ExchangeRate{
		&cringletest.ExchangeRate{From: "ABC", To: "DEF", Date: time.Now(), Value: decimal.New(1234, 3)},
	}

	err := sender.NotifyRates(context.Background(), rates)
	r.NoError(err)

	out, err := getTestOutput(buf)
	r.NoError(err)

	expected := fmt.Sprintf("Exchange Rate Results on %s:\n          1.0000    ABC Buys           1.2340    DEF\n", getFormattedDate(time.Now()))
	r.Equal(expected, out)
}

func TestNotifyRatesFailsWithNoRates(t *testing.T) {
	r := require.New(t)

	sender := New()

	err := sender.NotifyRates(context.Background(), nil)
	r.EqualError(err, cringletest.ErrNoRates.Error())
}

func TestNotifyValueSendsOK(t *testing.T) {
	r := require.New(t)

	sender, buf := getTestNotifier()

	rates := []*cringletest.ExchangeRate{
		&cringletest.ExchangeRate{From: "ABC", To: "DEF", Date: time.Now(), Value: decimal.New(1234, 3)},
	}

	err := sender.NotifyValue(context.Background(), decimal.New(2, 0), rates)
	r.NoError(err)

	out, err := getTestOutput(buf)
	r.NoError(err)

	expected := fmt.Sprintf("Currency Conversion Results on %s:\n          2.0000    ABC Buys           2.4680    DEF\n", getFormattedDate(time.Now()))
	r.Equal(expected, out)
}

func TestNotifyBestSendsOK(t *testing.T) {
	r := require.New(t)

	sender, buf := getTestNotifier()

	rate := &cringletest.ExchangeRate{From: "ABC", To: "DEF", Date: time.Now(), Value: decimal.New(1234, 3)}

	err := sender.NotifyBest(context.Background(), rate)
	r.NoError(err)

	out, err := getTestOutput(buf)
	r.NoError(err)

	expected := fmt.Sprintf("Best Exchange Rate in the last 7 days is:\n1.0000 ABC to 1.2340 DEF on %s\n", getFormattedDate(time.Now()))
	r.Equal(expected, out)
}
