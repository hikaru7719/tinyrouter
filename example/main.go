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
