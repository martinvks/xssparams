<h1 align="center">xssparams</h1>
<br>

`xssparams` takes a list of urls and identifies parameters potentially vulnerable to reflected xss

## Installation

```
go install github.com/martinvks/xssparams@latest
```

## Usage

For information about available flags, run:

```
xssparams -h
```

Example usage:

```
$ cat urls.txt
https://example.com?utm_source=google
https://example.com/articles/1
https://example.com/articles?query=computerphile
https://example.com?referer=https://youtube.com
$ cat urls.txt | xssparams
https://example.com/articles?query=computerphile [{query [SingleQuote]}]
https://example.com?referer=https://youtube.com [{referer [Href]}]
```

### What does the output mean?

- `Href` The parameter is reflected in the beggining of an href attribute
- `Element` The parameter is reflected inside an HTML element and the less-than sign is not escaped
- `Script` The parameter is reflected inside a script tag and the `</` character sequence is not escaped
- `DoubleQuote` The parameter is reflected inside double quotes and the double quote character is not escaped
- `SingleQuote` The parameter is reflected inside single quotes and the single quote character is not escaped or `\'` is escaped as `\\'`

### What is considered to be a parameter?

- Query Parameters, e.g., `search` and `language` in `https://example.com?search=quantum+computing&language=en`
- Numeric path segments, e.g., `123` in `https://example.com/articles/123`
- UUID path segments, e.g., `a92d7004-d18e-4aa3-9309-c016b6abca23` in `https://example.com/articles/a92d7004-d18e-4aa3-9309-c016b6abca23`

## Library Usage

```go
import (
	"github.com/martinvks/xssparams/pkg"
)

func main() {
	config := pkg.Config{
		Threads:      10,
		Timeout:      5,
		RateLimit:    50,
		CircuitBreak: 0,
		Verbose:      false,
		Headers:      map[string]string{"User-Agent": "custom-agent"},
		FilterCodes:  []int{404, 500},
	}

	urls := []string{
		"https://example.com?search=test",
		"https://example.com/articles/123",
	}

	results := pkg.Run(config, urls)

	for _, result := range results {
		// result.URL - the scanned URL
		// result.ParamsResults - slice of vulnerable parameters found
	}
}
```
