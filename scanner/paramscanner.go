package scanner

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/martinvks/xss-scanner/utils"
)

type ParamResult struct {
	param  string
	result []string
}

func scanParams(client *utils.RateLimitClient, target *url.URL, params []utils.Param) []ParamResult {
	var results []ParamResult

	for _, param := range params {
		paramResult, err := scanParam(client, target, param)
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

func scanParam(client *utils.RateLimitClient, target *url.URL, param utils.Param) ([]string, error) {
	id := utils.MiniUuid()

	resp, err := client.Get(getTargetUrl(target, param, id))
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
			resp, err := client.Get(getTargetUrl(target, param, id+input))
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
