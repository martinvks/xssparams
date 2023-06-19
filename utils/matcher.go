package utils

const (
	SingleQuote = "SingleQuote"
	DoubleQuote = "DoubleQuote"
	Element     = "Element"
	Href        = "Href"
	Unknown     = "Unknown"
)

type EscapeCheck struct {
	Checks    map[string]string
	MatchFunc Matcher
}

var EscapeChecks = map[string]EscapeCheck{
	SingleQuote: {
		Checks: map[string]string{
			`'`:  `'`,
			`\'`: `\\'`,
		},
		MatchFunc: MatchQuoteEnclosed,
	},
	DoubleQuote: {
		Checks: map[string]string{
			`"`: `"`,
		},
		MatchFunc: MatchDoubleQuoteEnclosed,
	},
	Element: {
		Checks: map[string]string{
			`<`: `<`,
		},
		MatchFunc: MatchBracketEnclosed,
	},
}

func FindMatchTypes(id string, body []byte) map[string]struct{} {
	matchTypes := make(map[string]struct{})

	if MatchQuoteEnclosed(id, body) {
		matchTypes[SingleQuote] = struct{}{}
	}
	if MatchDoubleQuoteEnclosed(id, body) {
		matchTypes[DoubleQuote] = struct{}{}
	}
	if MatchBracketEnclosed(id, body) {
		matchTypes[Element] = struct{}{}
	}
	if MatchHrefAttribute(id, body) {
		matchTypes[Href] = struct{}{}
	}

	if len(matchTypes) == 0 {
		if MatchAny(id, body) {
			matchTypes[Unknown] = struct{}{}
		}
	}

	return matchTypes
}
