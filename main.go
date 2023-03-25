package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/martinvks/xss-scanner/args"
	"github.com/martinvks/xss-scanner/scanner"
	"github.com/martinvks/xss-scanner/utils"
)

func main() {
	arguments, err := args.Parse()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	re := func(req *http.Request, via []*http.Request) error {
		differentHost := req.URL.Host != via[0].URL.Host

		if differentHost || len(via) >= 5 {
			return http.ErrUseLastResponse
		}
		return nil
	}

	client := &http.Client{
		CheckRedirect: re,
	}

	for _, val := range arguments.Urls {
		target, err := url.Parse(val)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to parse url: %s, %s\n", val, err)
			continue
		}

		params := utils.GetParams(target)
		if len(params) < 1 {
			continue
		}

		if arguments.Debug {
			fmt.Println(target.String())
		}

		results := scanner.ScanParams(client, target, params, arguments)
		if len(results) > 0 {
			fmt.Printf("%s: %v\n", target, results)
		}
	}
}
