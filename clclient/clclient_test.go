package clclient

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

func TestMain(m *testing.M) {
	loadDotEnv()

	_, err := envy.MustGet(cringletest.CurrencylayerAPIEnvVar)
	if err != nil {
		fmt.Println("Currencylayer API Key not found. Skipping clclient tests")
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestGetReturnsRequestedRates(t *testing.T) {
	r := require.New(t)

	client, err := New()
	r.NoError(err)

	rates, err := client.Get(context.Background(), "GBP", "EUR", "CAD")
	r.NoError(err)
	r.Contains(rates, "EUR")
	r.Contains(rates, "CAD")
}

func TestGetReturnsCorrectErrorWithBadAPIKey(t *testing.T) {
	r := require.New(t)

	var err error
	var client cringletest.RateClient

	envy.Temp(func() {
		envy.Set(cringletest.CurrencylayerAPIEnvVar, "NotAnAPIKey")
		client, err = New()
		r.NoError(err)
		_, err = client.Get(context.Background(), "GBP", "EUR", "CAD")
	})

	r.Error(err)
	r.EqualError(err, cringletest.ErrBadAuth.Error())
}

func TestGetReturnsCorrectValueWhenFromCurrencyDoesNotExist(t *testing.T) {
	r := require.New(t)

	client, err := New()
	r.NoError(err)

	_, err = client.Get(context.Background(), "NopeNopeNope", "EUR", "CAD")
	r.Error(err)
	r.EqualError(err, cringletest.ErrBadFromCurrency.Error())
}

func TestGetDoesNotReturnARateForUnknownCurrencies(t *testing.T) {
	r := require.New(t)

	client, err := New()
	r.NoError(err)

	rates, err := client.Get(context.Background(), "GBP", "EUR", "NopeNopeNope")
	r.NoError(err)
	r.Contains(rates, "EUR")
	r.NotContains(rates, "NopeNopeNope")
}

func TestGetReturnsCorrectErrorWhenNoToCurrenciesSupplied(t *testing.T) {
	r := require.New(t)

	client, err := New()
	r.NoError(err)

	_, err = client.Get(context.Background(), "GBP")
	r.Error(err)
	r.EqualError(err, cringletest.ErrNoToCurrencies.Error())
}

func TestGetOnReturnsRequestedRates(t *testing.T) {
	r := require.New(t)

	client, err := New()
	r.NoError(err)

	date, err := time.Parse("2006-01-02", "2016-05-12")
	r.NoError(err)

	rates, err := client.GetOn(context.Background(), date, "GBP", "EUR", "CAD")
	r.NoError(err)
	r.Contains(rates, "EUR")
	r.Contains(rates, "CAD")

	expectedVal, _ := new(decimal.Big).SetString("1.269817652318747")
	r.Equal(0, rates["EUR"].Value.Cmp(expectedVal))
}

func TestGetOnReturnsCorrectErrorWithBadAPIKey(t *testing.T) {
	r := require.New(t)

	var err error
	var client cringletest.RateClient

	envy.Temp(func() {
		envy.Set(cringletest.CurrencylayerAPIEnvVar, "NotAnAPIKey")
		client, err = New()
		r.NoError(err)
		_, err = client.GetOn(context.Background(), time.Now(), "GBP", "EUR", "CAD")
	})

	r.Error(err)
	r.EqualError(err, cringletest.ErrBadAuth.Error())
}

func TestGetOnReturnsCorrectValueWhenFromCurrencyDoesNotExist(t *testing.T) {
	r := require.New(t)

	client, err := New()
	r.NoError(err)

	_, err = client.GetOn(context.Background(), time.Now(), "NopeNopeNope", "EUR", "CAD")
	r.Error(err)
	r.EqualError(err, cringletest.ErrBadFromCurrency.Error())
}

func TestGetOnDoesNotReturnARateForUnknownCurrencies(t *testing.T) {
	r := require.New(t)

	client, err := New()
	r.NoError(err)

	rates, err := client.GetOn(context.Background(), time.Now(), "GBP", "EUR", "NopeNopeNope")
	r.NoError(err)
	r.Contains(rates, "EUR")
	r.NotContains(rates, "NopeNopeNope")
}

func TestGetOnReturnsCorrectErrorWhenNoToCurrenciesSupplied(t *testing.T) {
	r := require.New(t)

	client, err := New()
	r.NoError(err)

	_, err = client.GetOn(context.Background(), time.Now(), "GBP")
	r.Error(err)
	r.EqualError(err, cringletest.ErrNoToCurrencies.Error())
}
