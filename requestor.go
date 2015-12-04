// Package requestor makes calling RESTful API's easy.
package requestor

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Options contains parameters to be defined before sending the request to the server. Certain values
// can be omitted based on the request method (i.e. GET typically won't need to send a Body).
type Options struct {
	Method      string
	Query       map[string]string
	ContentType string
	Body        string
	Auth        []string
	Headers     map[string]string
}

// HTTPData contains the information returned from our request.
type HTTPData struct {
	Status  string
	Code    int
	Headers http.Header
	Body    []byte
	Error   error
}

// Send issues an HTTP request with the parameters specified in Options.
func Send(uri string, options *Options) *HTTPData {
	var req *http.Request
	var data HTTPData

	if options == nil {
		res, err := http.Get(uri)
		if err != nil {
			data.Error = err

			return &data
		}

		defer res.Body.Close()

		payload, _ := ioutil.ReadAll(res.Body)
		data.Body = payload
		data.Code = res.StatusCode
		data.Status = res.Status
		data.Headers = res.Header

		if res.StatusCode >= 400 {
			data.Error = fmt.Errorf("HTTP %d: %s", res.StatusCode, string(payload))
		}

		return &data
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	u, err := url.Parse(uri)
	if err != nil {
		data.Error = err

		return &data
	}

	query := u.Query()
	for k := range options.Query {
		query.Add(k, options.Query[k])
	}

	u.RawQuery = query.Encode()
	body := bytes.NewReader([]byte(options.Body))
	req, _ = http.NewRequest(strings.ToUpper(options.Method), u.String(), body)

	if len(options.Auth) > 0 {
		req.SetBasicAuth(options.Auth[0], options.Auth[1])
	}

	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	}

	if len(options.Headers) > 0 {
		for k := range options.Headers {
			req.Header.Add(k, options.Headers[k])
		}
	}

	res, err := client.Do(req)
	if err != nil {
		data.Error = err

		return &data
	}

	defer res.Body.Close()

	payload, _ := ioutil.ReadAll(res.Body)
	data.Body = payload
	data.Code = res.StatusCode
	data.Status = res.Status
	data.Headers = res.Header

	if res.StatusCode >= 400 {
		data.Error = fmt.Errorf("HTTP %d: %s", res.StatusCode, string(payload))
	}

	return &data
}
