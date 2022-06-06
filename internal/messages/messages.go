package messages

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/jalavosus/stuffnotifier/internal/logging"
	"github.com/jalavosus/stuffnotifier/internal/utils"
)

var logger = logging.NewLogger()

type Message interface {
	// FormatPlaintext should return the message text without any special formatting.
	FormatPlaintext() string
	// FormatMarkdown should return the message text with any desired markdown formatting.
	FormatMarkdown() string
	// PlaintextTemplate should return the template.Template (from the text/template package) used by Plaintext.
	PlaintextTemplate() *template.Template
	// MarkdownTemplate ahould return the template.Template (from the text/template package) used by Markdown.
	MarkdownTemplate() *template.Template
}

type baseMessage struct{}

func (m baseMessage) format(tmpl *template.Template, data any) (string, error) {
	b := new(bytes.Buffer)

	if err := tmpl.Execute(b, data); err != nil {
		err = errors.WithMessagef(err, "error executing template %[1]s", tmpl.Name())
		return "", err
	}

	return b.String(), nil
}

func formatTime(t time.Time) string {
	return utils.FormatTime(t, false)
}

func formatTimeOffset(t time.Time) string {
	return utils.FormatTime(t, true)
}

func formatDecimal(d decimal.Decimal) string {
	return d.String()
}

func formatCurrencyPair(a, b string) string {
	return fmt.Sprintf("%[1]s - %[2]s", a, b)
}

func formatPairQuote(a, b string, amtA, amtB decimal.Decimal) string {
	return fmt.Sprintf(
		"%[1]s %[2]s - %[3]s %[4]s",
		formatDecimal(amtA), a,
		formatDecimal(amtB), b,
	)
}

func formatTimezone(t time.Time, useLocalTime bool, timezone string) string {
	t = t.UTC()

	if useLocalTime {
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			log.Panicln(err)
		}

		t = t.In(loc)

		return formatTime(t)
	}

	return utils.FormatTimeUTC(t)
}

func isValidTime(t time.Time) bool {
	return !t.IsZero()
}
