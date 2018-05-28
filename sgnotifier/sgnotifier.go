// Package sgnotifier implemeents the cringletest.Notifier interface using the
// SendGrid api to send notifications via email
package sgnotifier

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ericlagergren/decimal"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
	"github.com/robotlovesyou/cringletest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type notifier struct {
	recipient string
	sender    string
	client    *sendgrid.Client
}

type formattedRate struct {
	From           string
	To             string
	OriginalValue  string
	ConvertedValue string
	Value          string
	Date           string
}

// Templates for outgoing email. If they grow larger then they should
// be placed into html templates and bundled into the binary using gobuffalo/packr or similar
const notifyRatesTemplate = `
<p><strong>Hello,</strong></p>
<p><strong>here are the rates you requested for <%= date %></strong></p>
<table>
	<%= for (rate) in rates { %>
		<tr>
			<td><%= rate.OriginalValue %> <%= rate.From %></td><td>Will buy you</td><td><%= rate.ConvertedValue %> <%= rate.To %></td>
		</tr>
	<% } %>
</table>
`

const notifyValuesTemplate = `
<p><strong>Hello,</strong></p>
<p><strong>here are the currency conversions requested for <%= date %></strong></p>
<table>
	<%= for (rate) in rates { %>
		<tr>
			<td><%= rate.OriginalValue %> <%= rate.From %></td><td>Will buy you</td><td><%= rate.ConvertedValue %> <%= rate.To %></td>
		</tr>
	<% } %>
</table>
`

const notifyBestTemplate = `
<p><strong>Hello,</strong><p>
<p>The best rate between <%= from %> and <%= to %> in the last 7 days was <%= rate %> on <%= date %></p>
`

const (
	ratesSubject  = "Your exchange rates"
	valuesSubject = "Your currency conversions"
	bestSubject   = "Your best exchange rate"
)

// New returns a new cringletest.Notifier which will send emails via sendgrid
func New(to string) (cringletest.Notifier, error) {
	apiKey, err := envy.MustGet(cringletest.SendGridAPIEnvVar)
	if err != nil {
		return nil, cringletest.ErrNoAuth
	}

	from, err := envy.MustGet(cringletest.SendGridFromAddressEnvVar)
	if err != nil {
		return nil, errors.New("no from address configured")
	}
	return &notifier{to, from, sendgrid.NewSendClient(apiKey)}, nil
}

func renderRateList(template string, rates []*formattedRate) (string, error) {
	ctx := plush.NewContext()
	ctx.Set("rates", rates)
	ctx.Set("date", rates[0].Date)

	s, err := plush.Render(template, ctx)
	if err != nil {
		return "", errors.Wrap(err, "could not render rate list")
	}

	return s, nil
}

func renderBest(template string, rate *formattedRate) (string, error) {
	ctx := plush.NewContext()
	ctx.Set("from", rate.From)
	ctx.Set("to", rate.To)
	ctx.Set("rate", rate.ConvertedValue)
	ctx.Set("date", rate.Date)

	s, err := plush.Render(template, ctx)
	if err != nil {
		return "", errors.Wrap(err, "could not render best")
	}

	return s, nil
}

func formatRate(originalValue *decimal.Big, rate *cringletest.ExchangeRate) *formattedRate {
	return &formattedRate{
		From:           rate.From,
		To:             rate.To,
		Date:           rate.Date.Format("Mon 02 Jan 2006"),
		OriginalValue:  fmt.Sprintf("%.4f", originalValue),
		ConvertedValue: fmt.Sprintf("%.4f", new(decimal.Big).Mul(originalValue, rate.Value)),
	}
}

func formatRates(originalValue *decimal.Big, rates []*cringletest.ExchangeRate) []*formattedRate {
	formattedRates := []*formattedRate{}
	for _, rate := range rates {
		formattedRates = append(formattedRates, formatRate(originalValue, rate))
	}

	return formattedRates
}

func (n *notifier) sendMail(subject, html string) error {
	from := mail.NewEmail("", n.sender)
	to := mail.NewEmail("", n.recipient)
	message := mail.NewSingleEmail(from, subject, to, html, html)

	response, err := n.client.Send(message)
	if err != nil {
		return cringletest.ErrSendFailed
	}
	if response.StatusCode == http.StatusUnauthorized {
		return cringletest.ErrBadAuth
	}
	if response.StatusCode != http.StatusAccepted {
		return cringletest.ErrSendFailed
	}
	return nil
}

func (n *notifier) notifyList(template, subject string, originalValue *decimal.Big, rates []*cringletest.ExchangeRate) error {
	if len(rates) == 0 {
		return cringletest.ErrNoRates
	}

	html, err := renderRateList(template, formatRates(originalValue, rates))
	if err != nil {
		return errors.Wrap(err, "could not notify rates")
	}

	return n.sendMail(subject, html)
}

func (n *notifier) NotifyRates(ctx context.Context, rates []*cringletest.ExchangeRate) error {
	return n.notifyList(notifyRatesTemplate, ratesSubject, decimal.New(1, 0), rates)
}

func (n *notifier) NotifyValue(ctx context.Context, value *decimal.Big, rates []*cringletest.ExchangeRate) error {
	return n.notifyList(notifyValuesTemplate, valuesSubject, value, rates)
}

func (n *notifier) NotifyBest(ctx context.Context, rate *cringletest.ExchangeRate) error {
	html, err := renderBest(notifyBestTemplate, formatRate(decimal.New(1, 0), rate))
	if err != nil {
		return errors.Wrap(err, "could not notify best")
	}

	return n.sendMail(bestSubject, html)
}
