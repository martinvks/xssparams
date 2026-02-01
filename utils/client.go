package utils

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"golang.org/x/time/rate"
)

type Response struct {
	Status  int
	Headers http.Header
	Body    []byte
}

type RateLimitClient struct {
	verbose           bool
	headers           map[string]string
	client            *http.Client
	limiter           *rate.Limiter
	circuitBreak      int
	consecutiveErrors int
	mu                sync.Mutex
}

func NewClient(headers map[string]string, rateLimit int, timeout int, circuitBreak int, verbose bool) *RateLimitClient {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	limiter := rate.NewLimiter(
		rate.Every(time.Second/time.Duration(rateLimit)),
		1,
	)

	return &RateLimitClient{
		verbose:      verbose,
		headers:      headers,
		client:       client,
		limiter:      limiter,
		circuitBreak: circuitBreak,
	}
}

func (c *RateLimitClient) CircuitBroken() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.circuitBreak > 0 && c.consecutiveErrors >= c.circuitBreak
}

func (c *RateLimitClient) Get(url string) (*Response, error) {
	req, err := newGetRequest(url, c.headers)
	if err != nil {
		return nil, err
	}

	err = c.limiter.Wait(req.Context())
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.circuitBreak > 0 {
		c.mu.Lock()
		if resp.StatusCode >= 400 {
			c.consecutiveErrors++
		} else {
			c.consecutiveErrors = 0
		}
		c.mu.Unlock()
	}

	if c.verbose {
		fmt.Printf("%s [%s]\n", req.URL, colorizedStatus(resp.StatusCode))
	}

	return &Response{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Body:    body,
	}, nil
}

func newGetRequest(url string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	req.Close = true

	return req, nil
}

func colorizedStatus(statusCode int) string {
	switch {
	case statusCode >= 400:
		return color.RedString("%d", statusCode)
	case statusCode >= 300:
		return color.YellowString("%d", statusCode)
	default:
		return color.GreenString("%d", statusCode)
	}
}
