// Copyright (c) 2021 Hikaru Miyahara
package tinyrouter

import (
	"net/http"
)

type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	Head(string, http.HandlerFunc)
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
	Put(string, http.HandlerFunc)
	Patch(string, http.HandlerFunc)
	Delete(string, http.HandlerFunc)
	Options(string, http.HandlerFunc)
	Connect(string, http.HandlerFunc)
	Trace(string, http.HandlerFunc)
}

// TinyRouter implements http.Handler interface.
type TinyRouter struct {
	routes []*route
}

func New() *TinyRouter {
	return &TinyRouter{
		routes: make([]*route, 0),
	}
}

func (t *TinyRouter) addRoute(method string, path string, f http.HandlerFunc) {
	route, err := newRoute(method, path, f)
	if err != nil {
		// TODO: implement err handling
		panic(err)
	}
	t.routes = append(t.routes, route)
}

func (t *TinyRouter) Head(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodHead, path, f)
}

func (t *TinyRouter) Get(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodGet, path, f)
}

func (t *TinyRouter) Post(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodPost, path, f)
}

func (t *TinyRouter) Put(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodPut, path, f)
}

func (t *TinyRouter) Patch(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodPatch, path, f)
}

func (t *TinyRouter) Delete(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodDelete, path, f)
}

func (t *TinyRouter) Options(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodOptions, path, f)
}

func (t *TinyRouter) Connect(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodConnect, path, f)
}

func (t *TinyRouter) Trace(path string, f http.HandlerFunc) {
	t.addRoute(http.MethodTrace, path, f)
}

func (t *TinyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler func(http.ResponseWriter, *http.Request)
	params := make(map[string]interface{})

	for _, route := range t.routes {
		isMatch := route.match(r)
		if isMatch {
			handler = route.handleFunc
			route.setParams(r, params)
			break
		}
	}

	if handler == nil {
		// TODO: set custom NotFoundHandler
		handler = defaultNotFoundHandler
	}
	newR := newRequest(r, params)
	handler(w, newR)
}
