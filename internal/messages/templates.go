package messages

import (
	"text/template"
)

var tmplFnMap = template.FuncMap{
	"FormatTime":       formatTime,
	"FormatTimeOffset": formatTimeOffset,
	"FormatDecimal":    formatDecimal,
	"FormatPair":       formatCurrencyPair,
	"FormatPairQuote":  formatPairQuote,
	"FormatTimezone":   formatTimezone,
	"IsValidTime":      isValidTime,
}

func mustParseTemplate(name, raw string) *template.Template {
	tmpl, tmplErr := template.New(name).
		Funcs(tmplFnMap).
		Parse(raw)

	tmpl = template.Must(tmpl, tmplErr)

	return tmpl
}

const (
	rawSpotPriceAlertPlaintextTemplate = `Crypto Spot Price Alert!

At {{ FormatTimeOffset .EventTime }}, {{ FormatPair .BaseCurrency .QuoteCurrency }} was {{ FormatPairQuote .BaseCurrency .QuoteCurrency 1.0 .SpotPrice }}.
`
	rawSpotPriceAlertMarkdownTemplate = ``
)

var (
	spotPriceAlertPlaintextTemplate = mustParseTemplate("spotPriceAlertPlaintext", rawSpotPriceAlertPlaintextTemplate)
	spotPriceAlertMarkdownTemplate  = mustParseTemplate("spotPriceAlertMarkdown", rawSpotPriceAlertMarkdownTemplate)
)

const (
	rawDeparturePlaintextTemplate = `--- Flight Information Update ---
{{ if .IsGateDeparture }}
Flight {{ .FlightNumber }} departed from {{ .Origin.Airport }} gate {{ .Origin.Gate }} at {{ FormatTimezone .GateDepartureTime.Actual .UseLocalTimezone .Origin.Timezone }}.
{{ if IsValidTime .TakeoffTime.Estimated }}
Estimated takeoff time is {{ FormatTimezone .TakeoffTime.Estimated .UseLocalTimezone .Origin.Timezone }}.
{{- end }}
{{ else if .IsTakeoff }}
Flight {{ .FlightNumber }} took off from {{ .Origin.Airport }} at {{ FormatTimezone .TakeoffTime.Actual .UseLocalTimezone .Origin.Timezone }}.
{{ if IsValidTime .LandingTime.Estimated }}
Estimated landing time at {{ .Destination.Airport }} is {{ FormatTimezone .LandingTime.Estimated .UseLocalTimezone .Destination.Timezone }} ({{ FormatTimezone .LandingTime.Estimated .UseLocalTimezone .Origin.Timezone }}).
{{- end -}}
{{- end -}}`

	rawArrivalPlaintextTemplate = `--- Flight Information Update ---
{{ if .IsGateArrival }}
{{- if .Destination.Terminal }}
Flight {{ .FlightNumber }} arrived at {{ .Destination.Airport }} (terminal {{ .Destination.Terminal }} gate {{ .Destination.Gate }}) at {{ FormatTimezone .GateArrivalTime.Actual .UseLocalTimezone .Destination.Timezone }} ({{ FormatTimezone .GateArrivalTime.Actual .UseLocalTimezone .Origin.Timezone }}).
{{ else }}
Flight {{ .FlightNumber }} arrived at {{ .Destination.Airport }} gate {{ .Destination.Gate }} at {{ FormatTimezone .GateArrivalTime.Actual .UseLocalTimezone .Destination.Timezone }} ({{ FormatTimezone .GateArrivalTime.Actual .UseLocalTimezone .Origin.Timezone }}).
{{- end -}}
{{ else if .IsLanding }}
Flight {{ .FlightNumber }} landed at {{ .Destination.Airport }} at {{ FormatTimezone .LandingTime.Actual .UseLocalTimezone .Destination.Timezone }} ({{ FormatTimezone .LandingTime.Actual .UseLocalTimezone .Origin.Timezone }}).
{{ if IsValidTime .GateArrivalTime.Estimated }}
Estimated gate arrival time is {{ FormatTimezone .GateArrivalTime.Estimated .UseLocalTimezone .Destination.Timezone }} ({{ FormatTimezone .GateArrivalTime.Estimated .UseLocalTimezone .Origin.Timezone }}).
{{- end -}}
{{- end -}}`

	rawPreDeparturePlaintextTemplate = `--- Flight Pre-Departure Alert ---

Flight {{ .FlightNumber }} is scheduled to depart from {{ .Origin.Airport }} at {{ FormatTimezone .GateDepartureTime.Scheduled .UseLocalTimezone .Origin.Timezone }}.
{{- if IsValidTime .GateDepartureTime.Estimated }}
{{ if .Origin.Terminal }}
Estimated departure time from Terminal {{ .Origin.Terminal }} Gate {{ .Origin.Gate }} is {{ FormatTimezone .GateDepartureTime.Estimated .UseLocalTimezone .Origin.Timezone }}.
{{ else }}
Estimated departure time from Gate {{ .Origin.Gate }} is {{ FormatTimezone .GateDepartureTime.Estimated .UseLocalTimezone .Origin.Timezone }}.
{{ end -}}
{{- end -}}`

	rawPreArrivalPlaintextTemplate = `--- Flight Pre-Arrival Alert ---

Flight {{ .FlightNumber }} is estimated to land at {{ .Destination.Airport }} at {{ FormatTimezone .LandingTime.Estimated .UseLocalTimezone .Destination.Timezone }} ({{ FormatTimezone .LandingTime.Estimated .UseLocalTimezone .Origin.Timezone }}).
{{- if IsValidTime .GateArrivalTime.Estimated }}
{{ if .Destination.Terminal }}
Estimated arrival at Terminal {{ .Destination.Terminal }} Gate {{ .Destination.Gate }} is {{ FormatTimezone .GateArrivalTime.Estimated .UseLocalTimezone .Destination.Timezone }} ({{ FormatTimezone .GateArrivalTime.Estimated .UseLocalTimezone .Origin.Timezone }}).
{{ else }}
Estimated arrival at Gate {{ .Destination.Gate }} is {{ FormatTimezone .GateArrivalTime.Estimated .UseLocalTimezone .Destination.Timezone }} ({{ FormatTimezone .GateArrivalTime.Estimated .UseLocalTimezone .Origin.Timezone }}).
{{ end -}}
{{- end -}}`
)

var (
	DepartureAlertPlaintextTemplate    = mustParseTemplate("departureAlertPlaintext", rawDeparturePlaintextTemplate)
	ArrivalAlertPlaintextTemplate      = mustParseTemplate("arrivalAlertPlaintext", rawArrivalPlaintextTemplate)
	PreDepartureAlertPlaintextTemplate = mustParseTemplate("preDeparturePlaintext", rawPreDeparturePlaintextTemplate)
	PreArrivalAlertPlaintextTemplate   = mustParseTemplate("preArrivalAlertPlaintext", rawPreArrivalPlaintextTemplate)
)
