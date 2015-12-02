package requestor

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Options builds our request parameters before sending it to the server.
type Options struct {
	Body        string
	ContentType string
	Auth        []string
}

// Get issues a HTTP GET request
func Get(url string, options *Options) ([]byte, error) {
	var req *http.Request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	body := bytes.NewReader([]byte(options.Body))
	req, _ = http.NewRequest("GET", url, body)

	if len(options.Auth) > 0 {
		req.SetBasicAuth(options.Auth[0], options.Auth[1])
	}

	if options.ContentType != "" {
		req.Header.Set("Content-Type", options.ContentType)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("HTTP %d :: %s", res.StatusCode, string(data[:])))
	}

	return data, nil
}
