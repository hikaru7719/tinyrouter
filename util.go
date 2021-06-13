package tinyrouter

import (
	"context"
	"net/http"
)

type contextKey string

var tinyrouterKey contextKey = "tinyrouterKey"

func DefaultNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Route Not Found"))
}

func NewRequest(r *http.Request, params map[string]interface{}) *http.Request {
	if len(params) == 0 {
		return r
	}
	ctx := context.WithValue(r.Context(), tinyrouterKey, params)
	return r.WithContext(ctx)
}

type Context struct {
	ctx context.Context
}

func NewContext(ctx context.Context) *Context {
	return &Context{ctx: ctx}
}

func Param(r *http.Request) *Context {
	ctx := r.Context()
	return NewContext(ctx)
}

func (c *Context) getString(key string) string {
	if m, ok := c.ctx.Value(tinyrouterKey).(map[string]interface{}); !ok {
		return ""
	} else if v, ok := m[key]; !ok {
		return ""
	} else if str, ok := v.(string); !ok {
		return ""
	} else {
		return str
	}
}

func (c *Context) getInt(key string) int {
	if m, ok := c.ctx.Value(tinyrouterKey).(map[string]interface{}); !ok {
		return 0
	} else if v, ok := m[key]; !ok {
		return 0
	} else if num, ok := v.(int); !ok {
		return 0
	} else {
		return num
	}
}
