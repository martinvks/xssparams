package main

import (
	"fmt"
	"github.com/martinvks/xss-scanner/args"
	"github.com/martinvks/xss-scanner/utils"
	"golang.org/x/exp/maps"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	MatchString  = "MATCH_STRING"
	MatchBracket = "MATCH_BRACKET"
	MatchOther   = "MATCH_OTHER"
)

type Response struct {
	status  int
	headers map[string][]string
	body    []byte
}

type Result []ParamResult

type ParamResult struct {
	param  string
	result map[string][]string
}

func main() {
	arguments, err := args.Parse()

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client := &http.Client{}

	for _, val := range arguments.Urls {
		req, err := utils.GetRequest(val, arguments.Headers)
		if err != nil {
			continue
		}

		resp, err := doRequest(client, req, arguments.Debug)
		if err != nil || resp.status != 200 {
			continue
		}

		target, err := url.Parse(val)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to parse url: %s, %s\n", val, err)
			continue
		}

		params := target.Query()
		if len(params) < 1 {
			_, _ = fmt.Fprintf(os.Stderr, "url has no query params: %s\n", target)
			continue
		}

		result := scanParams(client, cloneUrl(target), arguments)
		if len(result) > 0 {
			fmt.Printf("%s: %+v\n", target, result)
		}
	}
}

func scanParams(client *http.Client, target *url.URL, arguments args.Arguments) Result {
	var urlResult Result
	params := target.Query()

	for queryKey := range params {
		id := utils.MiniUuid()

		newParams := maps.Clone(params)
		newParams.Set(queryKey, id)
		target.RawQuery = newParams.Encode()

		req, err := utils.GetRequest(target.String(), arguments.Headers)
		if err != nil {
			continue
		}
		resp, err := doRequest(client, req, arguments.Debug)
		if err != nil {
			continue
		}

		results := make(map[string][]string)
		if utils.MatchStringEnclosed(id, resp.body) {
			results[MatchString] = []string{}
		}
		if utils.MatchBracketEnclosed(id, resp.body) {
			results[MatchBracket] = []string{}
		}
		if len(results) == 0 {
			if utils.MatchAny(id, resp.body) {
				results[MatchOther] = []string{}
			}
		}

		if len(results) == 0 {
			continue
		}

		for result := range results {
			var matchValue string
			var escapeChars []string
			var matchFunc utils.Matcher

			switch result {
			case MatchString:
				matchValue = id + "\""
				escapeChars = []string{"\"", "%22"}
				matchFunc = utils.MatchStringEnclosed
			case MatchBracket:
				matchValue = id + "<"
				escapeChars = []string{"<", "%3C"}
				matchFunc = utils.MatchBracketEnclosed
			default:
				continue
			}

			for _, escapeChar := range escapeChars {
				newParams.Set(queryKey, id+escapeChar)
				target.RawQuery = utils.EncodeExceptKey(newParams, queryKey)

				req, err := utils.GetRequest(target.String(), arguments.Headers)
				if err != nil {
					continue
				}
				resp, err := doRequest(client, req, arguments.Debug)
				if err != nil {
					continue
				}

				if matchFunc(matchValue, resp.body) {
					results[result] = append(results[result], escapeChar)
				}
			}
		}

		urlResult = append(urlResult, ParamResult{
			param:  queryKey,
			result: results,
		})
	}
	return urlResult
}

func doRequest(client *http.Client, req *http.Request, debug bool) (*Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Printf("%s: %d\n", req.URL, resp.StatusCode)
	}

	return &Response{
		status:  resp.StatusCode,
		headers: resp.Header,
		body:    body,
	}, nil
}

func cloneUrl(u *url.URL) *url.URL {
	u2 := *u
	return &u2
}
