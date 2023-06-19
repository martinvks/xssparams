package utils

import (
	"net/url"
	"testing"

	"golang.org/x/exp/slices"
)

type testCase struct {
	name     string
	url      string
	expected []Param
}

func TestGetParams(t *testing.T) {
	testCases := []testCase{
		{
			name:     "url with no params",
			expected: nil,
			url:      "https://example.com/comments",
		},
		{
			name: "url with query param",
			expected: []Param{{
				ParamType: QueryParam,
				ParamKey:  "query",
			}},
			url: "https://example.com/search?query=computerphile",
		},
		{
			name: "url with numeric path param",
			expected: []Param{{
				ParamType: PathParam,
				ParamKey:  "123",
				Index:     2,
			}},
			url: "https://example.com/comments/123",
		},
		{
			name: "url with UUID path param",
			expected: []Param{{
				ParamType: PathParam,
				ParamKey:  "a92d7004-d18e-4aa3-9309-c016b6abca23",
				Index:     2,
			}},
			url: "https://example.com/comments/a92d7004-d18e-4aa3-9309-c016b6abca23",
		},
		{
			name: "url with query and path param",
			expected: []Param{
				{
					ParamType: PathParam,
					ParamKey:  "425",
					Index:     1,
				},
				{
					ParamType: QueryParam,
					ParamKey:  "query",
				},
			},
			url: "https://example.com/425?query=quantum+computing",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			target, _ := url.Parse(tc.url)
			params := GetParams(target)

			if !slices.Equal(params, tc.expected) {
				t.Errorf("get params = %+v; want %+v", params, tc.expected)
			}
		})
	}
}
