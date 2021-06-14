package tinyrouter

import (
	"context"
	"net/http"
)

type contextKey string

var tinyrouterKey contextKey = "tinyrouterKey"

func defaultNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Route Not Found"))
}

func newRequest(r *http.Request, params map[string]interface{}) *http.Request {
	if len(params) == 0 {
		return r
	}
	ctx := context.WithValue(r.Context(), tinyrouterKey, params)
	return r.WithContext(ctx)
}

func Param(r *http.Request, key string) string {
	if m, ok := r.Context().Value(tinyrouterKey).(map[string]interface{}); !ok {
		return ""
	} else if v, ok := m[key]; !ok {
		return ""
	} else if str, ok := v.(string); !ok {
		return ""
	} else {
		return str
	}
}
