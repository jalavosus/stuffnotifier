package messages

import (
	"text/template"
	"time"

	"github.com/shopspring/decimal"
)

type SpotPriceAlert struct {
	baseMessage
	EventTime     time.Time
	SpotPrice     decimal.Decimal
	BaseCurrency  string
	QuoteCurrency string
}

func (a SpotPriceAlert) FormatPlaintext() string {
	msg, err := a.format(a.PlaintextTemplate(), a)
	if err != nil {
		panic(err)
	}

	return msg
}

func (a SpotPriceAlert) FormatMarkdown() string {
	msg, err := a.format(a.MarkdownTemplate(), a)
	if err != nil {
		panic(err)
	}

	return msg
}

func (a SpotPriceAlert) PlaintextTemplate() *template.Template {
	return spotPriceAlertPlaintextTemplate
}

func (a SpotPriceAlert) MarkdownTemplate() *template.Template {
	return spotPriceAlertMarkdownTemplate
}
