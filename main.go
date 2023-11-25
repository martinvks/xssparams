package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/martinvks/xssparams/args"
	"github.com/martinvks/xssparams/scanner"
	"github.com/martinvks/xssparams/utils"
)

func main() {
	arguments, err := args.Parse()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client := utils.NewClient(
		arguments.Headers,
		arguments.RateLimit,
		arguments.Timeout,
		arguments.Verbose,
	)

	scanWG := sync.WaitGroup{}
	scanChan := make(chan string)
	var results []*scanner.URLResult

	for i := 0; i < arguments.Threads; i++ {
		scanWG.Add(1)

		go func() {
			defer scanWG.Done()
			for url := range scanChan {
				result := scanner.Scan(client, url, arguments.FilterCodes)

				if result != nil {
					printURLResult(result)
					results = append(results, result)
				}
			}
		}()
	}

	for _, url := range arguments.Urls {
		scanChan <- url
	}
	close(scanChan)
	scanWG.Wait()

	if arguments.Verbose {
		printSummary(results)
	}
}

func printSummary(results []*scanner.URLResult) {
	fmt.Println()
	if len(results) == 0 {
		fmt.Println("Scan completed without any findings.")
		return
	}
	fmt.Printf("Found %d URLs potentially vulnerable to xss:\n", len(results))
	for _, result := range results {
		printURLResult(result)
	}
}

func printURLResult(result *scanner.URLResult) {
	fmt.Printf(
		"%s %s\n",
		result.URL,
		color.MagentaString("%v", result.ParamsResults),
	)
}
