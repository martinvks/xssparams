package scanner

import (
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/martinvks/xssparams/utils"
)

func Scan(client *utils.RateLimitClient, targetUrl string, filterCodes []int) []ParamResult {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return nil
	}

	params := utils.GetParams(target)
	if len(params) < 1 {
		return nil
	}

	resp, err := client.Get(target.String())
	if err != nil {
		return nil
	}

	contentType := resp.Headers.Get("content-type")
	if !strings.Contains(contentType, "html") && contentType != "" {
		return nil
	}

	if strings.HasPrefix(strconv.Itoa(resp.Status), "3") {
		return nil
	}

	if filterCodes != nil {
		if !slices.Contains(filterCodes, resp.Status) {
			return nil
		}
	}

	return scanParams(client, target, params)
}
