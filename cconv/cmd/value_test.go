package cmd

import (
	"context"
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
	"github.com/robotlovesyou/cringletest"
	"github.com/stretchr/testify/require"
)

func getFetchAndConvertArgs(cl cringletest.RateClient, notifiers []cringletest.Notifier) *requestConfig {
	return &requestConfig{
		From:      "ABC",
		To:        []string{"DEF", "GHI", "JKL"},
		Client:    cl,
		Notifiers: notifiers,
		Value:     decimal.New(2, 0),
	}
}

func TestFetchAndConvertNow(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(nil, nil, nil)

	err := fetchAndConvert(context.Background(), getFetchAndConvertArgs(cl, []cringletest.Notifier{n1, n2}))
	r.NoError(err)
}

func TestFetchAndConvertOnDate(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(nil, nil, nil)

	err := fetchAndConvert(context.Background(), getFetchAndConvertArgs(cl, []cringletest.Notifier{n1, n2}))
	r.NoError(err)
}

func TestFetchAndConvertReturnsCorrectClientError(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(cringletest.ErrBadAuth, nil, nil)

	err := fetchAndConvert(context.Background(), getFetchAndConvertArgs(cl, []cringletest.Notifier{n1, n2}))
	r.EqualError(errors.Cause(err), cringletest.ErrBadAuth.Error())
}

func TestFetchAndConvertReturnsCorrectNotifierError(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(nil, nil, cringletest.ErrBadAuth)

	err := fetchAndConvert(context.Background(), getFetchAndConvertArgs(cl, []cringletest.Notifier{n1, n2}))
	r.EqualError(errors.Cause(err), cringletest.ErrBadAuth.Error())
}
