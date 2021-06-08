package tinyrouter

import "net/http"

type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	Head(string, func(http.ResponseWriter, *http.Request))
	Get(string, func(http.ResponseWriter, *http.Request))
	Post(string, func(http.ResponseWriter, *http.Request))
	Put(string, func(http.ResponseWriter, *http.Request))
	Patch(string, func(http.ResponseWriter, *http.Request))
	Delete(string, func(http.ResponseWriter, *http.Request))
	Options(string, func(http.ResponseWriter, *http.Request))
	Connect(string, func(http.ResponseWriter, *http.Request))
	Trace(string, func(http.ResponseWriter, *http.Request))

	Use(func(http.ResponseWriter, *http.Request) func(http.ResponseWriter, *http.Request))
}
