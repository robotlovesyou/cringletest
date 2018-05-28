package sgnotifier

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gobuffalo/envy"
	"github.com/robotlovesyou/cringletest"
	"github.com/stretchr/testify/require"
)

func loadDotEnv() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory")
	}

	dotenvpath := path.Join(path.Join(wd, ".."), ".env")
	// .env may not exist so ignore any returned error from envy.Load
	envy.Load(dotenvpath)
}

func ensureEnvVars(keys ...string) error {
	for _, key := range keys {
		if _, err := envy.MustGet(key); err != nil {
			fmt.Printf("%s not found. Skipping sgnotifier tests\n", key)
			return err
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	loadDotEnv()

	if err := ensureEnvVars(
		cringletest.SendGridAPIEnvVar,
		cringletest.SendGridFromAddressEnvVar,
		cringletest.SendGridTestToAddressEnvVar); err != nil {

		os.Exit(0)
	}

	os.Exit(m.Run())
}

func getTestRecipient() string {
	recipient, err := envy.MustGet(cringletest.SendGridTestToAddressEnvVar)
	if err != nil {
		log.Fatalf("Could not get to address for tests")
	}

	return recipient
}

func TestNotifyRatesSendsOK(t *testing.T) {
	r := require.New(t)

	sender, err := New(getTestRecipient())
	r.NoError(err)

	rates := []*cringletest.ExchangeRate{
		&cringletest.ExchangeRate{From: "ABC", To: "DEF", Date: time.Now(), Value: decimal.New(1234, 3)},
	}

	err = sender.NotifyRates(context.Background(), rates)
	r.NoError(err)
}

func TestNotifyRatesFailsWithBadAPIKey(t *testing.T) {
	r := require.New(t)

	var err error
	var sender cringletest.Notifier

	rates := []*cringletest.ExchangeRate{
		&cringletest.ExchangeRate{From: "ABC", To: "DEF", Date: time.Now(), Value: decimal.New(1234, 3)},
	}

	envy.Temp(func() {
		envy.Set(cringletest.SendGridAPIEnvVar, "Not an api key")
		sender, err = New(getTestRecipient())
		r.NoError(err)

		err = sender.NotifyRates(context.Background(), rates)
		r.EqualError(err, cringletest.ErrBadAuth.Error())
	})
}

func TestNotifyRatesFailsWithNoRates(t *testing.T) {
	r := require.New(t)

	sender, err := New(getTestRecipient())
	r.NoError(err)

	err = sender.NotifyRates(context.Background(), nil)
	r.EqualError(err, cringletest.ErrNoRates.Error())
}

func TestNotifyValueSendsOK(t *testing.T) {
	r := require.New(t)

	sender, err := New(getTestRecipient())
	r.NoError(err)

	rates := []*cringletest.ExchangeRate{
		&cringletest.ExchangeRate{From: "ABC", To: "DEF", Date: time.Now(), Value: decimal.New(1234, 3)},
	}

	err = sender.NotifyValue(context.Background(), decimal.New(2, 0), rates)
	r.NoError(err)
}

func TestNotifyBestSendsOK(t *testing.T) {
	r := require.New(t)

	sender, err := New(getTestRecipient())
	r.NoError(err)

	rate := &cringletest.ExchangeRate{From: "ABC", To: "DEF", Date: time.Now(), Value: decimal.New(1234, 3)}

	err = sender.NotifyBest(context.Background(), rate)
	r.NoError(err)
}
