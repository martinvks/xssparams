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

var regexNumeric = regexp.MustCompile(`^\d+$`)
var regexUuid = regexp.MustCompile(`^(?i)[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}$`)

func GetParams(target *url.URL) []Param {
	var params []Param

	segments := strings.Split(target.Path, "/")
	for index, segment := range segments {
		if regexNumeric.MatchString(segment) || regexUuid.MatchString(segment) {
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
