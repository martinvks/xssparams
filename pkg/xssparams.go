package pkg

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
	"github.com/martinvks/xssparams/scanner"
	"github.com/martinvks/xssparams/utils"
)

type Config struct {
	Threads      int
	Timeout      int
	RateLimit    int
	CircuitBreak int
	Verbose      bool
	Headers      map[string]string
	FilterCodes  []int
}

func Run(config Config, urls []string) []*scanner.URLResult {
	client := utils.NewClient(
		config.Headers,
		config.RateLimit,
		config.Timeout,
		config.CircuitBreak,
		config.Verbose,
	)

	scanWG := sync.WaitGroup{}
	scanChan := make(chan string)
	var results []*scanner.URLResult
	var mu sync.Mutex

	for i := 0; i < config.Threads; i++ {
		scanWG.Add(1)

		go func() {
			defer scanWG.Done()
			for url := range scanChan {
				result := scanner.Scan(client, url, config.FilterCodes)

				if result != nil {
					printURLResult(result)
					mu.Lock()
					results = append(results, result)
					mu.Unlock()
				}
			}
		}()
	}

	for _, url := range urls {
		if client.CircuitBroken() {
			break
		}
		scanChan <- url
	}
	close(scanChan)
	scanWG.Wait()

	return results
}

func printURLResult(result *scanner.URLResult) {
	fmt.Printf(
		"%s %s\n",
		result.URL,
		color.MagentaString("%v", result.ParamsResults),
	)
}
