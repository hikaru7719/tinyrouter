// Copyright (c) 2021 Hikaru Miyahara
package tinyrouter

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTinyRouter(t *testing.T) {
	cases := []struct {
		method       string
		path         string
		expectStatus int
		expectBody   string
	}{
		{
			method:       http.MethodGet,
			path:         "/todo",
			expectStatus: 200,
			expectBody:   "GET /todo",
		},
		{
			method:       http.MethodPost,
			path:         "/todo",
			expectStatus: 200,
			expectBody:   "POST /todo",
		},
		{
			method:       http.MethodPut,
			path:         "/todo/abc",
			expectStatus: 200,
			expectBody:   "PUT /todo/abc",
		},
		{
			method:       http.MethodDelete,
			path:         "/todo/abc",
			expectStatus: 200,
			expectBody:   "DELETE /todo/abc",
		},
		{
			method:       http.MethodGet,
			path:         "/todo/abc/def",
			expectStatus: 200,
			expectBody:   "GET /todo/abc/def",
		},
		{
			method:       http.MethodGet,
			path:         "/notfound",
			expectStatus: 404,
			expectBody:   "Route Not Found",
		},
	}

	router := New()
	testSetupRouter(t, router)

	server := httptest.NewServer(router)
	defer server.Close()

	for _, tc := range cases {
		tc := tc
		req, _ := http.NewRequest(tc.method, server.URL+tc.path, nil)
		res, _ := http.DefaultClient.Do(req)
		b, _ := io.ReadAll(res.Body)
		assert.Equal(t, tc.expectBody, string(b))
	}
}

func testSetupRouter(t *testing.T, router *TinyRouter) {
	routes := []struct {
		method string
		path   string
		f      func(w http.ResponseWriter, r *http.Request)
	}{
		{
			method: http.MethodGet,
			path:   "/todo",
			f: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("GET /todo"))
			},
		},
		{
			method: http.MethodPost,
			path:   "/todo",
			f: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("POST /todo"))
			},
		},
		{
			method: http.MethodPut,
			path:   "/todo/{id}",
			f: func(w http.ResponseWriter, r *http.Request) {
				id := Param(r, "id")
				w.Write([]byte(fmt.Sprintf("PUT /todo/%s", id)))
			},
		},
		{
			method: http.MethodDelete,
			path:   "/todo/{id}",
			f: func(w http.ResponseWriter, r *http.Request) {
				id := Param(r, "id")
				w.Write([]byte(fmt.Sprintf("DELETE /todo/%s", id)))
			},
		},
		{
			method: http.MethodGet,
			path:   "/todo/{id}/{field}",
			f: func(w http.ResponseWriter, r *http.Request) {
				id := Param(r, "id")
				field := Param(r, "field")
				w.Write([]byte(fmt.Sprintf("GET /todo/%s/%s", id, field)))
			},
		},
	}

	for _, route := range routes {
		testSetRoute(t, router, route.method, route.path, route.f)
	}
}

func testSetRoute(t *testing.T, router *TinyRouter, method string, path string, f func(w http.ResponseWriter, r *http.Request)) {
	switch method {
	case http.MethodGet:
		router.Get(path, f)
	case http.MethodPost:
		router.Post(path, f)
	case http.MethodPut:
		router.Put(path, f)
	case http.MethodDelete:
		router.Delete(path, f)
	default:
		t.Fatal("unexpected http method")
	}
}
