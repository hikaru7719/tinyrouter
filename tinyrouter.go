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

// TinyRouter implements http.Handler interface.
type TinyRouter struct{}

func NewRouter() *TinyRouter {
	return &TinyRouter{}
}

func (t *TinyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (t *TinyRouter) Head(path string, f func(http.ResponseWriter, *http.Request)) {}

func (t *TinyRouter) Get(path string, f func(http.ResponseWriter, *http.Request)) {}

func (t *TinyRouter) Post(path string, f func(http.ResponseWriter, *http.Request)) {}

func (t *TinyRouter) Put(path string, f func(http.ResponseWriter, *http.Request)) {}

func (t *TinyRouter) Patch(path string, f func(http.ResponseWriter, *http.Request)) {}

func (t *TinyRouter) Delete(path string, f func(http.ResponseWriter, *http.Request)) {}

func (t *TinyRouter) Options(path string, f func(http.ResponseWriter, *http.Request)) {}

func (t *TinyRouter) Connect(path string, f func(http.ResponseWriter, *http.Request)) {}

func (t *TinyRouter) Trace(path string, f func(http.ResponseWriter, *http.Request)) {}
