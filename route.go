package tinyrouter

import (
	"net/http"
)

type Route struct {
	Method     string
	Path       string
	HandleFunc func(http.ResponseWriter, *http.Request)
}

func NewRoute(method string, path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return &Route{
		Method:     method,
		Path:       path,
		HandleFunc: f,
	}
}

func (m *Route) Match(r *http.Request) bool {
	return false
}
