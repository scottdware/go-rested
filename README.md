## go-requestor
[![GoDoc](https://godoc.org/github.com/scottdware/go-requestor?status.svg)](https://godoc.org/github.com/scottdware/go-requestor) [![Travis-CI](https://travis-ci.org/scottdware/go-requestor.svg?branch=master)](https://travis-ci.org/scottdware/go-requestor)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/scottdware/go-requestor/master/LICENSE)

A Go package that makes calling RESTful API's easy.

### Example

To send a request, you first specify some options such as the method (GET, POST, PUT, etc.), Content-Type, any additional headers, etc.:

```Go
req := &requestor.Options{
	Method:      "get",
	Auth:        []string{"user", "securepassword"},
	ContentType: "application/json",
}
```

Then, call the `Send()` function to issue the request, specifying the URL and the options you defined above:

```Go
data, err := requestor.Send("https://someurl/api/v1.0/stuff", req)
if err != nil {
	fmt.Println(err)
}
```

The returned data is a struct with the following values:

```Go
type HTTPData struct {
	Body    []byte
	Status  string
	Code    int
	Headers http.Header
	Error   error
}
```

For instance, you can see the payload by converting the `Body` field to a string:

```Go
fmt.Println(string(data.Body))
```