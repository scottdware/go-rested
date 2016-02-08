// Package rested makes calling RESTful API's easy.
package rested

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Request contains parameters to be defined before sending the request to the server. Certain values
// can be omitted based on the request method (i.e. GET typically won't need to send a Body).
type Request struct {
	Method      string
	Query       map[string]string
	ContentType string
	Body        string
	Auth        []string
	Headers     map[string]string
}

// Response contains the information returned from our request.
type Response struct {
	Status  string
	Code    int
	Headers http.Header
	Body    []byte
	Error   error
}

// Send issues an HTTP request with the parameters specified in request.
func Send(uri string, request *Request) *Response {
	var req *http.Request
	var data Response

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	if request == nil {
		res, err := client.Get(uri)
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

	u, err := url.Parse(uri)
	if err != nil {
		data.Error = err

		return &data
	}

	query := u.Query()
	for k := range request.Query {
		query.Add(k, request.Query[k])
	}

	u.RawQuery = query.Encode()
	body := bytes.NewReader([]byte(request.Body))
	req, _ = http.NewRequest(strings.ToUpper(request.Method), u.String(), body)

	if len(request.Auth) > 0 {
		req.SetBasicAuth(request.Auth[0], request.Auth[1])
	}

	if len(request.ContentType) > 0 {
		req.Header.Set("Content-Type", request.ContentType)
	}

	if len(request.Headers) > 0 {
		for k := range request.Headers {
			req.Header.Add(k, request.Headers[k])
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
