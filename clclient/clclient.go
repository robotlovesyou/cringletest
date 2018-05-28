// Package clclient implements the cringletest/RateClient interface using the
// currencylayer API
package clclient

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/robotlovesyou/cringletest"
	"gopkg.in/resty.v1"
)

const (
	currencyLayerEndpoint = "http://apilayer.net/api/"
	liveMethod            = "live"
	historicalMethod      = "historical"
	accessKeyName         = "access_key"
	currenciesName        = "currencies"
	invalidAccessKey      = "invalid_access_key"
	dateName              = "date"
)

type rateClient struct {
	apiKey string
}

// Wrap decimal.Big in a struct with a custom JSON unmarshal func to allow it to be unmarshalled from json
// because currencylayer returns its decimal rates as floating point numbers. Genius! /s
type jsonBig struct {
	*decimal.Big
}

func (jb *jsonBig) UnmarshalJSON(b []byte) error {
	jb.Big = new(decimal.Big)
	return jb.UnmarshalText(b)
}

type clError struct {
	Code int64  `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

type clResult struct {
	Success bool                `json:"success"`
	Error   *clError            `json:"error,omitempty"`
	Date    string              `json:"date"`
	Source  string              `json:"source"`
	Quotes  map[string]*jsonBig `json:"quotes"`
}

// New returns a new RateClient
func New() (cringletest.RateClient, error) {
	apiKey, err := envy.MustGet(cringletest.CurrencylayerAPIEnvVar)
	if err != nil {
		return nil, cringletest.ErrNoAuth
	}

	return &rateClient{apiKey}, nil
}

func parse(from string, result *clResult) (cringletest.RateMap, error) {
	// if the result.Success is not true then return an appropriate error
	// if the result.Success is true then iterate through the returned currencies and create the
	// ExchangeRate results
	if !result.Success {
		if result.Error.Type == invalidAccessKey {
			return nil, cringletest.ErrBadAuth
		}

		return nil, fmt.Errorf("currencylayer api error %s", result.Error.Type)
	}

	var date time.Time
	var err error
	if len(result.Date) == 0 {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", result.Date)
		if err != nil {
			return nil, errors.Wrap(err, "bad date returned by currencylayer api")
		}
	}

	fromRate, ok := result.Quotes[result.Source+from]
	if !ok {
		return nil, cringletest.ErrBadFromCurrency
	}
	delete(result.Quotes, result.Source+from)

	rates := cringletest.RateMap{}

	for name, rate := range result.Quotes {
		er := &cringletest.ExchangeRate{
			From:  from,
			To:    name[len(result.Source):],
			Date:  date,
			Value: new(decimal.Big).Quo(rate.Big, fromRate.Big),
		}
		rates[er.To] = er
	}

	return rates, nil
}

func (rc *rateClient) apiRequest(ctx context.Context, method string, params map[string]string) (*clResult, error) {
	result := new(clResult)
	params[accessKeyName] = rc.apiKey

	// if the context already has a deadline set dont set a new one, otherwise use the CurrencyLayerTimeout
	_, ok := ctx.Deadline()
	var cancel context.CancelFunc
	if !ok {
		ctx, cancel = context.WithTimeout(context.Background(), cringletest.CurrencylayerTimeout)
		defer cancel()
	}

	resp, err := resty.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(result).
		Get(fmt.Sprintf("%s%s", currencyLayerEndpoint, method))

	if err != nil {
		return nil, errors.Wrap(err, "could not make currencylayer api request")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("unexpected status code returned from currencylayer api")
	}

	return result, nil
}

// Implements cringletest.RateClient.Get using the currencylayer api.
func (rc *rateClient) Get(ctx context.Context, from string, to ...string) (cringletest.RateMap, error) {
	if len(to) == 0 {
		return nil, cringletest.ErrNoToCurrencies
	}

	currencies := append(to, from)
	params := map[string]string{
		currenciesName: strings.Join(currencies, ","),
	}

	result, err := rc.apiRequest(ctx, liveMethod, params)
	if err != nil {
		return nil, errors.Wrap(err, "could not get live currencies")
	}
	return parse(from, result)
}

// Implements cringletest.RateClient.GetOn using the currencylayer api.
func (rc *rateClient) GetOn(ctx context.Context, date time.Time, from string, to ...string) (cringletest.RateMap, error) {
	if len(to) == 0 {
		return nil, cringletest.ErrNoToCurrencies
	}

	currencies := append(to, from)
	params := map[string]string{
		dateName:       date.Format("2006-01-02"),
		currenciesName: strings.Join(currencies, ","),
	}

	result, err := rc.apiRequest(ctx, historicalMethod, params)
	if err != nil {
		return nil, errors.Wrap(err, "could not get live currencies")
	}
	return parse(from, result)
}
