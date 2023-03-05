package utils

import "net/http"

func GetRequest(url string, headers map[string]string) (*http.Request, error) {
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
