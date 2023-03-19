package scanners

import (
	"fmt"
	"strings"
)

type ParamResult struct {
	param  string
	result map[string]struct{}
}

func (p ParamResult) String() string {
	keys := make([]string, len(p.result))

	i := 0
	for k := range p.result {
		keys[i] = k
		i++
	}

	return fmt.Sprintf("{%s: %s}", p.param, strings.Join(keys, ", "))
}
