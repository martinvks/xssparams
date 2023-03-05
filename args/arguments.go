package args

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Arguments struct {
	Debug   bool
	Headers map[string]string
	Urls    []string
}

type headersFlag []string

func (i *headersFlag) String() string {
	return ""
}

func (i *headersFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	debug        bool
	headersSlice headersFlag
	headersMap   = make(map[string]string)
	urls         []string
)

func Parse() (Arguments, error) {
	flag.BoolVar(&debug, "debug", false, "log request urls and status codes")
	flag.Var(&headersSlice, "H", "header fields added to each request. syntax similar to curl: -H \"x-header: val\".")
	flag.Parse()

	for _, header := range headersSlice {
		name, value, found := strings.Cut(header, ":")
		if !found {
			return Arguments{}, fmt.Errorf("invalid header '%s', expected syntax: 'x-header: val'", header)
		}

		trimmedName := strings.TrimSpace(name)
		trimmedValue := strings.TrimSpace(value)
		headersMap[trimmedName] = trimmedValue
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return Arguments{}, fmt.Errorf("failed to read input: %s\n", err)
	}

	return Arguments{
		Debug:   debug,
		Headers: headersMap,
		Urls:    urls,
	}, nil
}
