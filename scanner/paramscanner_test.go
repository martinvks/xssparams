package scanner

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/martinvks/xss-scanner/utils"
)

type testCase struct {
	name               string
	expected           []ParamResult
	responseDataWriter func(string) string
}

func TestParamScanner(t *testing.T) {
	testCases := []testCase{
		{
			name:     "escaped value reflected inside HTML element",
			expected: nil,
			responseDataWriter: func(s string) string {
				return fmt.Sprintf(
					"<html><body><p>Search value: %s</p></body></html>",
					html.EscapeString(s),
				)
			},
		},
		{
			name:     "escaped value reflected inside HTML attribute",
			expected: nil,
			responseDataWriter: func(s string) string {
				return fmt.Sprintf(
					"<html><body><input type=\"text\" value=\"%s\"></body></html>",
					html.EscapeString(s),
				)
			},
		},
		{
			name: "value reflected inside HTML element",
			expected: []ParamResult{{
				"q",
				[]string{"Element"},
			}},
			responseDataWriter: func(s string) string {
				return fmt.Sprintf(
					"<html><body><p>Search value: %s</p></body></html>",
					s,
				)
			},
		},
		{
			name: "value reflected inside HTML attribute",
			expected: []ParamResult{{
				"q",
				[]string{"DoubleQuote"},
			}},
			responseDataWriter: func(s string) string {
				return fmt.Sprintf(
					"<html><body><input type=\"text\" value=\"%s\"></body></html>",
					s,
				)
			},
		},
		{
			name: "value reflected inside single quote HTML attribute",
			expected: []ParamResult{{
				"q",
				[]string{"SingleQuote"}},
			},
			responseDataWriter: func(s string) string {
				return fmt.Sprintf(
					"<html><body><input type='text' value='%s'></body></html>",
					s,
				)
			},
		},
		{
			name: "value reflected inside script tag",
			expected: []ParamResult{{
				"q",
				[]string{"Script"}},
			},
			responseDataWriter: func(s string) string {
				return fmt.Sprintf(
					"<html><body><script type=\"application/ld+json\">{\"query\": \"%s\"}</script></body></html>",
					strings.ReplaceAll(s, `"`, `\"`),
				)
			},
		},
		{
			name: "value reflected in beginning of href attribute",
			expected: []ParamResult{{
				"q",
				[]string{"Href"}},
			},
			responseDataWriter: func(s string) string {
				return fmt.Sprintf(
					"<html><body><a href=\"%s\">Click me</a></body></html>",
					html.EscapeString(s),
				)
			},
		},
		{
			name: "value reflected inside HTML element and attribute",
			expected: []ParamResult{{
				"q",
				[]string{"DoubleQuote", "Element"}},
			},
			responseDataWriter: func(s string) string {
				return fmt.Sprintf(
					"<html><body><input type=\"text\" value=\"%[1]s\"><p>Search value: %[1]s</p></body></html>",
					s,
				)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				searchParam := r.URL.Query().Get("q")
				response := tc.responseDataWriter(searchParam)
				_, _ = fmt.Fprintf(w, response)
			}))
			defer ts.Close()

			target, _ := url.Parse(ts.URL + "/search?q=computerphile")
			params := []utils.Param{{
				ParamKey:  "q",
				ParamType: utils.QueryParam,
			}}

			paramResults := scanParams(utils.NewClient(nil, 100, false), target, params)

			if len(paramResults) != len(tc.expected) {
				t.Errorf("len(paramResults) = %d; want %d", len(paramResults), len(tc.expected))
			}

			for index := range paramResults {
				result := paramResults[index]
				expected := tc.expected[index]

				if result.param != expected.param {
					t.Errorf("param name = \"%s\"; want \"%s\"", result.param, expected.param)

				}

				if !slices.Equal(result.result, expected.result) {
					t.Errorf("param result = %s; want %s", result.result, expected.result)
				}
			}
		})
	}
}
