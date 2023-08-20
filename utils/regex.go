package utils

import "regexp"

type Matcher func(string, []byte) bool

func MatchAny(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(escaped)
	return matcher.Match(body)
}

func MatchSingleQuoteContext(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`'[^<>'"]*` + escaped + `[^<>'"]*'`)
	return matcher.Match(body)
}

func MatchDoubleQuoteContext(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`"[^<>'"]*` + escaped + `[^<>'"]*"`)
	return matcher.Match(body)
}

func MatchElementContext(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`>[^<>'"]*` + escaped + `[^"<>]*<`)
	return matcher.Match(body)
}

func MatchScriptContext(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`<script[^>]*>([^<]|<[^/])*` + escaped + `([^<]|<[^/])*</script`)
	return matcher.Match(body)
}

func MatchHrefAttribute(value string, body []byte) bool {
	escaped := regexp.QuoteMeta(value)
	matcher := regexp.MustCompile(`href="` + escaped)
	return matcher.Match(body)
}
