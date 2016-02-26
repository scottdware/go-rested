// Package rested makes calling RESTful API's easy.
package rested

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	// "strconv"
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

var (
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

// NewRequest creates the state for our REST call.
func NewRequest() *Request {
	return &Request{}
}

// BasicAuth sets the authentication using a standard username/password combination.
func (r *Request) BasicAuth(user, password string) {
	r.Auth = []string{user, password}
}

// Send issues an HTTP request with the given options. To send/post form values, use a map[string]string
// type for the body parameter.
func (r *Request) Send(method, uri string, body interface{}, headers, query map[string]string) *Response {
	var req *http.Request
	var data Response

	u, err := url.Parse(uri)
	if err != nil {
		data.Error = err

		return &data
	}

	q := u.Query()

	if query != nil {
		for k := range query {
			q.Add(k, query[k])
		}
	}

	u.RawQuery = q.Encode()

	switch body.(type) {
	case []byte:
		b := bytes.NewReader(body.([]byte))
		req, _ = http.NewRequest(strings.ToUpper(method), u.String(), b)
	case map[string]string:
		form := url.Values{}
		for k, v := range body.(map[string]string) {
			form.Add(k, v)
		}

		req, _ = http.NewRequest(strings.ToUpper(method), u.String(), strings.NewReader(form.Encode()))
		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if len(r.Auth) > 0 {
		req.SetBasicAuth(r.Auth[0], r.Auth[1])
	}

	if headers != nil {
		for k := range headers {
			req.Header.Add(k, headers[k])
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
