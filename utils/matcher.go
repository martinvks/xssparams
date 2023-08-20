package utils

const (
	SingleQuote = "SingleQuote"
	DoubleQuote = "DoubleQuote"
	Element     = "Element"
	Script      = "Script"
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
		MatchFunc: MatchSingleQuoteContext,
	},
	DoubleQuote: {
		Checks: map[string]string{
			`"`: `"`,
		},
		MatchFunc: MatchDoubleQuoteContext,
	},
	Element: {
		Checks: map[string]string{
			`<`: `<`,
		},
		MatchFunc: MatchElementContext,
	},
	Script: {
		Checks: map[string]string{
			`</`: `</`,
		},
		MatchFunc: MatchScriptContext,
	},
}

func FindMatchTypes(id string, body []byte) map[string]struct{} {
	matchTypes := make(map[string]struct{})

	if MatchSingleQuoteContext(id, body) {
		matchTypes[SingleQuote] = struct{}{}
	}
	if MatchDoubleQuoteContext(id, body) {
		matchTypes[DoubleQuote] = struct{}{}
	}
	if MatchElementContext(id, body) {
		matchTypes[Element] = struct{}{}
	}
	if MatchScriptContext(id, body) {
		matchTypes[Script] = struct{}{}
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
