package utils

import (
	"net/url"
	"sort"
	"strings"
)

func EncodeExceptKey(values url.Values, exceptKey string) string {
	if values == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := values[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')

			var value string
			if k == exceptKey {
				value = v
			} else {
				value = url.QueryEscape(v)
			}
			buf.WriteString(value)
		}
	}
	return buf.String()
}
