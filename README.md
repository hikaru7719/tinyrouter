# TinyRouter

TinyRouter is HTTP routing library to study routing algorithm.
This library is so simple and small. But this is not for using at production.

TinyRouter can only do path base routing.

# Install

```
go get -u github.com/hikaru7719/tinyrouter
```

# Example

```go
package main

import (
	"net/http"

	"github.com/hikaru7719/tinyrouter"
)

func main() {
	r := tinyrouter.New()
	r.Get("/hello", Hello)
	r.Get("/hello/{name}", HelloName)

	http.ListenAndServe(":8080", r)
}

func Hello(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello World!"))
}

func HelloName(rw http.ResponseWriter, r *http.Request) {
	name := tinyrouter.Param(r, "name")
	rw.Write([]byte("Hello " + name + "!"))
}
```
