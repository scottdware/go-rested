package requestor

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Options contains parameters to be defined before sending the request to the server. Certain values
// can be omitted based on the request method (i.e. GET typically won't need to send a Body).
type Options struct {
	Method      string
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
	Payload []byte
	Error   error
}

// Send issues an HTTP request with the values specified in Options.
func Send(url string, options *Options) *HTTPData {
	var req *http.Request
	var data HTTPData
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
		data.Error = err

		return &data
	}

	defer res.Body.Close()

	payload, _ := ioutil.ReadAll(res.Body)

	data.Payload = payload
	data.Code = res.StatusCode
	data.Status = res.Status
	data.Headers = res.Header

	if res.StatusCode >= 400 {
		data.Error = fmt.Errorf("HTTP %d: %s", res.StatusCode, string(payload))
	}

	return &data
}

// String will convert the payload/body of the request from a []byte to a string value.
func (h *HTTPData) String() string {
	return string(h.Payload)
}
