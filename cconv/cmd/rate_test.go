package cmd

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/robotlovesyou/cringletest"
	"github.com/robotlovesyou/cringletest/testclient"
	"github.com/robotlovesyou/cringletest/testnotifier"
	"github.com/stretchr/testify/require"
)

func getClientAndNotifiers(clientErr, notifier1Err, notifier2Err error) (cringletest.RateClient, cringletest.Notifier, cringletest.Notifier) {
	return testclient.New(clientErr), testnotifier.New(notifier1Err), testnotifier.New(notifier2Err)
}

func getFetchAndShowArgs(cl cringletest.RateClient, notifiers []cringletest.Notifier) *requestConfig {
	return &requestConfig{
		From:      "ABC",
		To:        []string{"DEF", "GHI", "JKL"},
		Client:    cl,
		Notifiers: notifiers,
	}
}

func TestFetchAndShowNow(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(nil, nil, nil)

	err := fetchAndShow(context.Background(), getFetchAndShowArgs(cl, []cringletest.Notifier{n1, n2}))
	r.NoError(err)
}

func TestFetchAndShowOnDate(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(nil, nil, nil)

	err := fetchAndShow(context.Background(), getFetchAndShowArgs(cl, []cringletest.Notifier{n1, n2}))
	r.NoError(err)
}

func TestFetchAndShowReturnsCorrectClientError(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(cringletest.ErrBadAuth, nil, nil)

	err := fetchAndShow(context.Background(), getFetchAndShowArgs(cl, []cringletest.Notifier{n1, n2}))
	r.EqualError(errors.Cause(err), cringletest.ErrBadAuth.Error())
}

func TestFetchAndShowReturnsCorrectNotifierError(t *testing.T) {
	r := require.New(t)
	cl, n1, n2 := getClientAndNotifiers(nil, nil, cringletest.ErrBadAuth)

	err := fetchAndShow(context.Background(), getFetchAndShowArgs(cl, []cringletest.Notifier{n1, n2}))
	r.EqualError(errors.Cause(err), cringletest.ErrBadAuth.Error())
}
