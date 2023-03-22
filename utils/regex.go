package utils

import "regexp"

type Matcher func(string, []byte) bool

func MatchAny(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(escaped)
	return matcher.Match(body)
}

func MatchQuoteEnclosed(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`'[^<>'"]*` + escaped + `[^<>'"]*'`)
	return matcher.Match(body)
}

func MatchDoubleQuoteEnclosed(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`"[^<>'"]*` + escaped + `[^<>'"]*"`)
	return matcher.Match(body)
}

func MatchBracketEnclosed(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`>[^<>'"]*` + escaped + `[^"<>]*<`)
	return matcher.Match(body)
}

func MatchHrefAttribute(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`href="` + escaped)
	return matcher.Match(body)
}
