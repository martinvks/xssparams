package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
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

	client := utils.NewClient(
		arguments.Headers,
		arguments.RateLimit,
		arguments.Verbose,
	)

	scanWG := sync.WaitGroup{}
	scanChan := make(chan string)

	for i := 0; i < arguments.Threads; i++ {
		scanWG.Add(1)

		go func() {
			defer scanWG.Done()
			for url := range scanChan {
				result := scanner.Scan(client, url, arguments.FilterCodes)

				if result != nil && len(result) > 0 {
					fmt.Printf("%s %s\n", url, color.MagentaString("%v", result))
				}
			}
		}()
	}

	for _, url := range arguments.Urls {
		scanChan <- url
	}
	close(scanChan)

	scanWG.Wait()
}
