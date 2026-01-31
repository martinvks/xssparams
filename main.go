package main

import (
	"fmt"
	"os"

	"github.com/martinvks/xssparams/args"
	"github.com/martinvks/xssparams/pkg"
	"github.com/martinvks/xssparams/scanner"
)

func main() {
	arguments, err := args.Parse()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	config := pkg.Config{
		Threads:      arguments.Threads,
		Timeout:      arguments.Timeout,
		RateLimit:    arguments.RateLimit,
		CircuitBreak: arguments.CircuitBreak,
		Verbose:      arguments.Verbose,
		Headers:      arguments.Headers,
		FilterCodes:  arguments.FilterCodes,
	}

	results := pkg.Run(config, arguments.Urls)
	printSummary(results)
}

func printSummary(results []*scanner.URLResult) {
	if len(results) == 0 {
		fmt.Println("Scan completed without any findings.")
		return
	}
	fmt.Printf("Found %d URLs potentially vulnerable to xss\n", len(results))
}
