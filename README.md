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
	"fmt"
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
	fmt.Fprint(rw, "Hello World!\n")
}

func HelloName(rw http.ResponseWriter, r *http.Request) {
	name := tinyrouter.Param(r, "name")
	fmt.Fprintf(rw, "Hello %s!\n", name)
}
```
