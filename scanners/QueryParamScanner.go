package scanners

import (
	"golang.org/x/exp/maps"
	"net/http"
	"net/url"

	"github.com/martinvks/xss-scanner/args"
	"github.com/martinvks/xss-scanner/utils"
)

func ScanQueryParams(client *http.Client, target *url.URL, arguments args.Arguments) []ParamResult {
	var results []ParamResult
	params := target.Query()

	if len(params) < 1 {
		return nil
	}

	for queryKey := range params {
		id := utils.MiniUuid()

		newParams := maps.Clone(params)
		newParams.Set(queryKey, id)
		target.RawQuery = newParams.Encode()

		resp, err := utils.DoRequest(client, target.String(), arguments)
		if err != nil {
			continue
		}

		matchTypes := utils.FindMatchTypes(id, resp.Body, resp.Headers)
		if len(matchTypes) == 0 {
			continue
		}

		var queryResults = make(map[string]struct{})

		for matchType := range matchTypes {
			matchCheck, ok := utils.MatchChecks[matchType]
			if !ok {
				queryResults[matchType] = struct{}{}
				continue
			}

			for _, input := range matchCheck.Inputs {
				newParams.Set(queryKey, id+input)
				target.RawQuery = utils.EncodeExceptKey(newParams, queryKey)

				resp, err := utils.DoRequest(client, target.String(), arguments)
				if err != nil {
					continue
				}

				if matchCheck.MatchFunc(id+matchCheck.Char, resp.Body) {
					queryResults[matchType] = struct{}{}
				}
			}
		}

		if len(queryResults) > 0 {
			results = append(results, ParamResult{
				param:  queryKey,
				result: queryResults,
			})
		}
	}
	return results
}
