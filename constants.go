package cringletest

import (
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrBadAuth should be returned when an authentication or authorization error occurs getting
	ErrBadAuth = errors.New("client authentication error")
	// ErrNoAuth should be returned when no authentication information is provided to a client
	ErrNoAuth = errors.New("client has no authentication")
)

const (
	// CurrencylayerAPIEnvVar is env var name for the currencylayer api key
	CurrencylayerAPIEnvVar = "CURRENCYLAYER_API_KEY"
	// CurrencylayerTimeout is the maximum time to wait for a call to the Currencylayer api
	CurrencylayerTimeout = 10 * time.Second
	// SendGridAPIEnvVar is the name of the SendGrid API Key Env Var
	SendGridAPIEnvVar = "SENDGRID_API_KEY"
	//SendGridFromAddressEnvVar is the env var containing the test from address for sendgrid tests
	SendGridFromAddressEnvVar = "SENDGRID_FROM_ADDRESS"
	//SendGridTestToAddressEnvVar is the env var containing the to address for sendgrid tests
	SendGridTestToAddressEnvVar = "SENDGRID_TEST_TO_ADDRESS"
	// NotifyTimeout is the timeout allocated to the notify functions
	NotifyTimeout = 30 * time.Second
)
