package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/martinvks/xss-scanner/args"
	"github.com/martinvks/xss-scanner/scanners"
	"github.com/martinvks/xss-scanner/utils"
)

func main() {
	arguments, err := args.Parse()

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client := &http.Client{}

	for _, val := range arguments.Urls {

		resp, err := utils.DoRequest(client, val, arguments)
		if err != nil || resp.Status != 200 {
			continue
		}

		target, err := url.Parse(val)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to parse url: %s, %s\n", val, err)
			continue
		}

		var results []scanners.ParamResult

		pathResults := scanners.ScanPathParams(client, cloneUrl(target), arguments)
		if len(pathResults) > 0 {
			results = append(results, pathResults...)
		}

		queryResults := scanners.ScanQueryParams(client, cloneUrl(target), arguments)
		if len(queryResults) > 0 {
			results = append(results, queryResults...)
		}

		if len(results) > 0 {
			fmt.Printf("%s: %v\n", target, results)
		}
	}
}

func cloneUrl(u *url.URL) *url.URL {
	u2 := *u
	return &u2
}
