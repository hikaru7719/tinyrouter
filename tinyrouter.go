package tinyrouter

import (
	"net/http"
)

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

// TinyRouter implements http.Handler interface.
type TinyRouter struct {
	Routes []*Route
}

func NewRouter() *TinyRouter {
	return &TinyRouter{
		Routes: make([]*Route, 0),
	}
}

func (t *TinyRouter) addRoute(method string, path string, f func(http.ResponseWriter, *http.Request)) {
	route, err := NewRoute(method, path, f)
	if err != nil {
		// TODO: implement err handling
		panic(err)
	}
	t.Routes = append(t.Routes, route)
}

func (t *TinyRouter) Head(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodHead, path, f)
}

func (t *TinyRouter) Get(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodGet, path, f)
}

func (t *TinyRouter) Post(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodPost, path, f)
}

func (t *TinyRouter) Put(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodPut, path, f)
}

func (t *TinyRouter) Patch(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodPatch, path, f)
}

func (t *TinyRouter) Delete(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodDelete, path, f)
}

func (t *TinyRouter) Options(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodOptions, path, f)
}

func (t *TinyRouter) Connect(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodConnect, path, f)
}

func (t *TinyRouter) Trace(path string, f func(http.ResponseWriter, *http.Request)) {
	t.addRoute(http.MethodTrace, path, f)
}

func (t *TinyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler func(w http.ResponseWriter, r *http.Request)
	params := make(map[string]interface{})

	for _, route := range t.Routes {
		isMatch := route.Match(r)
		if isMatch {
			handler = route.HandleFunc
			route.SetParams(r, params)
			break
		}
	}

	if handler == nil {
		// TODO: set custom NotFoundHandler
		handler = DefaultNotFoundHandler
	}
	newR := NewRequest(r, params)
	handler(w, newR)
}
