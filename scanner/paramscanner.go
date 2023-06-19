package scanner

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/martinvks/xss-scanner/args"
	"github.com/martinvks/xss-scanner/utils"
)

type ParamResult struct {
	param  string
	result []string
}

func ScanParams(client *http.Client, target *url.URL, params []utils.Param, arguments args.Arguments) []ParamResult {
	var results []ParamResult

	for _, param := range params {
		paramResult, err := scanParam(client, target, param, arguments)
		if err != nil || len(paramResult) < 1 {
			continue
		}

		results = append(results, ParamResult{
			param:  param.ParamKey,
			result: paramResult,
		})
	}
	return results
}

func scanParam(client *http.Client, target *url.URL, param utils.Param, arguments args.Arguments) ([]string, error) {
	id := utils.MiniUuid()

	resp, err := utils.DoRequest(client, getTargetUrl(target, param, id), arguments)
	if err != nil {
		return nil, err
	}

	var results []string
	matchTypes := utils.FindMatchTypes(id, resp.Body)

	for matchType := range matchTypes {
		escapeCheck, ok := utils.EscapeChecks[matchType]
		if !ok {
			results = append(results, matchType)
			continue
		}

		for input, match := range escapeCheck.Checks {
			resp, err := utils.DoRequest(client, getTargetUrl(target, param, id+input), arguments)
			if err != nil {
				continue
			}

			if escapeCheck.MatchFunc(id+match, resp.Body) {
				results = append(results, matchType)
				break
			}
		}
	}

	return results, nil
}

func getTargetUrl(target *url.URL, param utils.Param, newValue string) string {
	switch param.ParamType {
	case utils.PathParam:
		return getPathParamTargetUrl(cloneUrl(target), param.Index, newValue)
	case utils.QueryParam:
		return getQueryParamTargetUrl(cloneUrl(target), param.ParamKey, newValue)
	default:
		panic(fmt.Sprintf("unknown param type: %d", param.ParamType))
	}
}

func getPathParamTargetUrl(target *url.URL, index int, newValue string) string {
	segments := strings.Split(target.Path, "/")
	segments[index] = newValue
	target.Path = strings.Join(segments, "/")
	return target.String()
}

func getQueryParamTargetUrl(target *url.URL, key string, newValue string) string {
	params := target.Query()
	params.Set(key, newValue)
	target.RawQuery = params.Encode()
	return target.String()
}

func cloneUrl(u *url.URL) *url.URL {
	u2 := *u
	return &u2
}
