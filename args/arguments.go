package args

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var validStatusCode = regexp.MustCompile(`^[1-5]\d{2}$`)

type Arguments struct {
	Threads     int
	RateLimit   int
	Verbose     bool
	Headers     map[string]string
	FilterCodes []int
	Urls        []string
}

type headersFlag []string

func (i *headersFlag) String() string {
	return ""
}

func (i *headersFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type filterCodesFlag []int

func (i *filterCodesFlag) String() string {
	return ""
}

func (i *filterCodesFlag) Set(value string) error {
	if value == "" {
		return nil
	}

	statuseCodes := strings.Split(value, ",")

	for _, c := range statuseCodes {
		statusCode := strings.TrimSpace(c)

		if !validStatusCode.MatchString(statusCode) {
			return fmt.Errorf("invalid status code '%s'", statusCode)
		}

		code, err := strconv.Atoi(statusCode)
		if err != nil {
			return fmt.Errorf("failed to convert status code '%s' to int", statusCode)
		}

		*i = append(*i, code)
	}

	return nil
}

var (
	threads      int
	rateLimit    int
	verbose      bool
	filterCodes  filterCodesFlag
	headersSlice headersFlag
)

func Parse() (Arguments, error) {
	flag.IntVar(
		&threads,
		"threads",
		10,
		"number of lightweight threads to use",
	)
	flag.IntVar(
		&rateLimit,
		"rate-limit",
		50,
		"maximum requests to send per second",
	)
	flag.BoolVar(
		&verbose,
		"verbose",
		false,
		"print request urls and status codes",
	)
	flag.Var(
		&filterCodes,
		"filter-codes",
		"only scan urls that initially return one of the supplied status codes, e.g., -status-codes \"200,404\"",
	)
	flag.Var(
		&headersSlice,
		"H",
		"header fields added to each request. syntax similar to curl: -H \"x-header: val\"",
	)
	flag.Parse()

	if threads < 1 {
		return Arguments{}, fmt.Errorf("invalid value '%d' for flag -threads, must be greater than 0.", threads)
	}

	if rateLimit < 1 {
		return Arguments{}, fmt.Errorf("invalid value '%d' for flag -rate-limit, must be greater than 0.", rateLimit)
	}

	headers := make(map[string]string)
	for _, header := range headersSlice {
		name, value, found := strings.Cut(header, ":")
		if !found {
			return Arguments{}, fmt.Errorf("invalid header '%s' for flag -H, expected syntax: 'x-header: val'", header)
		}
		headers[strings.TrimSpace(name)] = strings.TrimSpace(value)
	}

	var urls []string
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return Arguments{}, fmt.Errorf("failed to read input: %s", err)
	}

	return Arguments{
		Threads:     threads,
		RateLimit:   rateLimit,
		Verbose:     verbose,
		Headers:     headers,
		FilterCodes: filterCodes,
		Urls:        urls,
	}, nil
}
