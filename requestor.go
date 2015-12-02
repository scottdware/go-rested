package requestor

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Options builds our request parameters before sending it to the server.
type Options struct {
	Method      string
	Body        string
	ContentType string
	Auth        []string
	Headers     map[string]string
}

// Send issues a HTTP request with the values specified in Options.
func Send(url string, options *Options) ([]byte, error) {
	var req *http.Request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	body := bytes.NewReader([]byte(options.Body))
	req, _ = http.NewRequest(strings.ToUpper(options.Method), url, body)

	if len(options.Auth) > 0 {
		req.SetBasicAuth(options.Auth[0], options.Auth[1])
	}

	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	}

	if len(options.Headers) > 0 {
		for k, _ := range options.Headers {
			req.Header.Add(k, options.Headers[k])
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("HTTP %d: %s", res.StatusCode, string(data[:])))
	}

	return data, nil
}
