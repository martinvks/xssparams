package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/martinvks/xss-scanner/args"
)

type Response struct {
	Status  int
	Headers http.Header
	Body    []byte
}

func DoRequest(client *http.Client, url string, arguments args.Arguments) (*Response, error) {
	req, err := getRequest(url, arguments.Headers)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
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

	if arguments.Debug {
		fmt.Printf("%s: %d\n", req.URL, resp.StatusCode)
	}

	return &Response{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Body:    body,
	}, nil
}

func getRequest(url string, headers map[string]string) (*http.Request, error) {
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
