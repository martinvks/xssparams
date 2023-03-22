package utils

import (
	"net/url"
	"regexp"
	"strings"
)

const (
	PathParam = iota
	QueryParam
)

type Param struct {
	ParamType int
	ParamKey  string
	Index     int
}

func GetParams(target *url.URL) []Param {
	var params []Param

	re := regexp.MustCompile(`^\d+$`)
	segments := strings.Split(target.Path, "/")
	for index, segment := range segments {
		if re.MatchString(segment) {
			params = append(params, Param{
				ParamType: PathParam,
				ParamKey:  segment,
				Index:     index,
			})
		}
	}

	queryParams := target.Query()
	for queryKey := range queryParams {
		params = append(params, Param{
			ParamType: QueryParam,
			ParamKey:  queryKey,
		})
	}

	return params
}
