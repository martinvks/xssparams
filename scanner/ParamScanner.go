package scanner

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/martinvks/xss-scanner/args"
	"github.com/martinvks/xss-scanner/utils"
)

func ScanParams(client *http.Client, target *url.URL, params []utils.Param, arguments args.Arguments) []ParamResult {
	var results []ParamResult

	for _, param := range params {
		id := utils.MiniUuid()

		resp, err := utils.DoRequest(client, getTargetUrl(target, param, id), arguments)
		if err != nil {
			continue
		}

		matchTypes := utils.FindMatchTypes(id, resp.Body, resp.Headers)
		if len(matchTypes) == 0 {
			continue
		}

		var paramResults = make(map[string]struct{})

		for matchType := range matchTypes {
			escapeCheck, ok := utils.EscapeChecks[matchType]
			if !ok {
				paramResults[matchType] = struct{}{}
				continue
			}

			for input, match := range escapeCheck.Checks {
				resp, err := utils.DoRequest(client, getTargetUrl(target, param, id+input), arguments)
				if err != nil {
					continue
				}

				if escapeCheck.MatchFunc(id+match, resp.Body) {
					paramResults[matchType] = struct{}{}
				}
			}
		}

		if len(paramResults) > 0 {
			results = append(results, ParamResult{
				param:  param.ParamKey,
				result: paramResults,
			})
		}

	}
	return results
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
