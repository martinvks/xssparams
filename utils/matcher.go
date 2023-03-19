package utils

import (
	"net/http"
	"strings"
)

const (
	MatchQuote       = "MATCH_QUOTE"
	MatchDoubleQuote = "MATCH_DOUBLE_QUOTE"
	MatchBracket     = "MATCH_BRACKET"
	MatchHref        = "MATCH_HREF"
	MatchOther       = "MATCH_OTHER"
	MatchHeader      = "MATCH_HEADER"
)

type MatchCheck struct {
	Char      string
	Inputs    []string
	MatchFunc Matcher
}

var MatchChecks = map[string]MatchCheck{
	MatchQuote: {
		Char:      "'",
		Inputs:    []string{"'", "%27"},
		MatchFunc: MatchQuoteEnclosed,
	},
	MatchDoubleQuote: {
		Char:      "\"",
		Inputs:    []string{"\"", "%22"},
		MatchFunc: MatchDoubleQuoteEnclosed,
	},
	MatchBracket: {
		Char:      "<",
		Inputs:    []string{"<", "%3C"},
		MatchFunc: MatchBracketEnclosed,
	},
}

func FindMatchTypes(id string, body []byte, headers http.Header) map[string]struct{} {
	matchTypes := make(map[string]struct{})

	if MatchQuoteEnclosed(id, body) {
		matchTypes[MatchQuote] = struct{}{}
	}
	if MatchDoubleQuoteEnclosed(id, body) {
		matchTypes[MatchDoubleQuote] = struct{}{}
	}
	if MatchBracketEnclosed(id, body) {
		matchTypes[MatchBracket] = struct{}{}
	}
	if MatchHrefAttribute(id, body) {
		matchTypes[MatchHref] = struct{}{}
	}

	if len(matchTypes) == 0 {
		if MatchAny(id, body) {
			matchTypes[MatchOther] = struct{}{}
		}
	}

	for _, headerValues := range headers {
		for _, headerValue := range headerValues {
			if strings.Contains(headerValue, id) {
				matchTypes[MatchHeader] = struct{}{}
			}
		}
	}
	return matchTypes
}
