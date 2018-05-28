package cmd

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/robotlovesyou/cringletest"
	"github.com/stretchr/testify/require"
)

func getBestRateArgs(cl cringletest.RateClient, notifiers []cringletest.Notifier) *requestConfig {
	return &requestConfig{
		From:      "ABC",
		To:        []string{"DEF"},
		Client:    cl,
		Notifiers: notifiers,
	}
}

func TestFetchBest(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(nil, nil, nil)

	err := fetchBest(context.Background(), getBestRateArgs(cl, []cringletest.Notifier{n1, n2}))
	r.NoError(err)
}

func TestFetchBestReturnsCorrectClientError(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(cringletest.ErrBadAuth, nil, nil)

	err := fetchBest(context.Background(), getBestRateArgs(cl, []cringletest.Notifier{n1, n2}))
	r.EqualError(errors.Cause(err), cringletest.ErrBadAuth.Error())
}

func TestFetchBestReturnsCorrectNotifierError(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(nil, nil, cringletest.ErrBadAuth)

	err := fetchBest(context.Background(), getBestRateArgs(cl, []cringletest.Notifier{n1, n2}))
	r.EqualError(errors.Cause(err), cringletest.ErrBadAuth.Error())
}
