package requestor

import (
	"bytes"
	"crypto/tls"
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
	if options.Auth != "" {
		req.SetBasicAuth(options.Auth[0], options.Auth[1])
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		// if res.Header["Content-Type"][0] == "application/json" {
		// }

		return data, errors.New(fmt.Sprintf("HTTP %d :: %s", res.StatusCode, string(data[:])))
	}

	return data, nil
}
