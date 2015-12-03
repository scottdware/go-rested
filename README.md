## go-requestor
[![GoDoc](https://godoc.org/github.com/scottdware/go-requestor?status.svg)](https://godoc.org/github.com/scottdware/go-requestor) [![Travis-CI](https://travis-ci.org/scottdware/go-requestor.svg?branch=master)](https://travis-ci.org/scottdware/go-requestor)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/scottdware/go-requestor/master/LICENSE)

A Go package that makes calling RESTful API's easy.

### Example

To send a request, you first specify some options such as the method (GET, POST, PUT, etc.), Content-Type, any additional headers, query parameters, etc.:

```Go
req := &requestor.Options{
	Method:      "get",
	ContentType: "application/json",
	Auth:        []string{"user", "securepassword"},
	Headers: map[string]string{
		"Custom-Header": "value",
	},
	Query: map[string]string{
		"another_param": "somevalue",
		"more_stuff": "more-values",
	},
}
```

Then, call the `Send()` function to issue the request, specifying the URL and the options you defined above:

```Go
data := requestor.Send("https://someurl/api/v1.0/stuff?default_param=something", req)
if data.Error != nil {
	fmt.Println(data.Error)
}
```
> If there was any type of error in your request, it will be defined in the `Error` field of the returned struct. This is why the above error check looks a bit different than how you would normally check for errors in Go.

The entire request with any additional query parameters defined will look like the following when sent to the server:

```
https://someurl/api/v1.0/stuff?another_param=somevalue&default_param=something&more_stuff=more-values
```

The returned data is a struct with the following values:

```Go
type HTTPData struct {
	Status  string
	Code    int
	Headers http.Header
	Body    []byte
	Error   error
}
```

If you would like to do a plain-and-simple GET request, without any authentication or parameters, you can leave the `Options` parameter `nil`, like so:

```Go
data := requestor.Send("https://someurl/api/v1.0/stuff", nil)
```

You can see the payload by converting the `Body` field to a string:

```Go
fmt.Println(string(data.Body))
```