package scanners

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/martinvks/xss-scanner/args"
	"github.com/martinvks/xss-scanner/utils"
)

func ScanPathParams(client *http.Client, target *url.URL, arguments args.Arguments) []ParamResult {
	var results []ParamResult
	path := target.Path
	pathParams := utils.FindPathParams(path)

	if len(pathParams) < 1 {
		return nil
	}

	for _, pathParam := range pathParams {
		id := utils.MiniUuid()
		target.Path = strings.Replace(path, pathParam, id, 1)

		resp, err := utils.DoRequest(client, target.String(), arguments)
		if err != nil {
			continue
		}

		matchTypes := utils.FindMatchTypes(id, resp.Body, resp.Headers)
		if len(matchTypes) == 0 {
			continue
		}

		var paramResults = make(map[string]struct{})

		for matchType := range matchTypes {
			matchCheck, ok := utils.MatchChecks[matchType]
			if !ok {
				paramResults[matchType] = struct{}{}
				continue
			}

			for _, input := range matchCheck.Inputs {
				target.Path = strings.Replace(path, pathParam, id+input, 1)

				resp, err := utils.DoRequest(client, target.String(), arguments)
				if err != nil {
					continue
				}

				if matchCheck.MatchFunc(id+matchCheck.Char, resp.Body) {
					paramResults[matchType] = struct{}{}
				}
			}
		}

		if len(paramResults) > 0 {
			results = append(results, ParamResult{
				param:  pathParam,
				result: paramResults,
			})
		}
	}
	return results
}
