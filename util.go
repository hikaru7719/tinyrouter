package tinyrouter

import "net/http"

var DefaultNotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Route Not Found"))
}
