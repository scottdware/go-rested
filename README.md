## go-rested
[![GoDoc](https://godoc.org/github.com/scottdware/go-rested?status.svg)](https://godoc.org/github.com/scottdware/go-rested) [![Travis-CI](https://travis-ci.org/scottdware/go-rested.svg?branch=master)](https://travis-ci.org/scottdware/go-rested)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/scottdware/go-rested/master/LICENSE)

A Go package that makes calling RESTful API's easy.

### Usage & Examples

Initiate the request, and set authentication (if applicable):

```Go
r := rested.NewRequest()
r.BasicAuth("user", "secret")
```

If you need to specify any headers or query parameters, then use a `map[string]string` type to define them:

```Go
headers := map[string]string{
	"Content-Type": "application/json",
	"Accept": "application/json",
}

query := map[string]string{
	"search_string": "dog",
	"results": "10",
}
```

Then, call the `Send()` function to issue the request, specifying the URL and the options you defined above:

```Go
data := r.Send("get", "https://someurl/api/v1.0/stuff?default_param=something", nil, headers, query)
```

> If there was any type of error in your request, it will be defined in the `Error` field of the returned struct. You can check for errors similar to how you normally do in Go:
```Go
if data.Error != nil {
	fmt.Println(data.Error)
}
```

The entire request with any additional query parameters defined will look like the following when sent to the server:

```
https://someurl/api/v1.0/stuff?results=10&search_string=dog&default_param=something
```

The returned data is a struct with the following fields:

```Go
type HTTPData struct {
	Status  string
	Code    int
	Headers http.Header
	Body    []byte
	Error   error
}
```

You can see the payload by converting the `Body` field to a string:

```Go
fmt.Println(string(data.Body))
```