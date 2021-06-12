package tinyrouter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractParam(t *testing.T) {
	cases := map[string]struct {
		path   string
		expect []string
	}{
		"path without param": {
			path:   "/todo",
			expect: []string{},
		},
		"path with param": {
			path:   "/todo/{id}",
			expect: []string{"id"},
		},
		"path with multiple param": {
			path:   "/todo/{id}/{field}",
			expect: []string{"id", "field"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			indexList, _ := bracesIndex(tc.path)
			actual := extractParam(tc.path, indexList)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestMakePattern(t *testing.T) {
	cases := map[string]struct {
		path   string
		expect string
	}{
		"path without param": {
			path:   "/todo",
			expect: "^/todo[/]?$",
		},
		"path with param": {
			path:   "/todo/{id}",
			expect: "^/todo/([^/]+)[/]?$",
		},
		"path with multiple param": {
			path:   "/todo/{id}/{field}",
			expect: "^/todo/([^/]+)/([^/]+)[/]?$",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			indexList, _ := bracesIndex(tc.path)
			actual := makePatternString(tc.path, indexList)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestNormalize(t *testing.T) {
	cases := map[string]struct {
		path        string
		expect      string
		expectError error
	}{
		"invalid path": {
			path:        "todo",
			expect:      "",
			expectError: InvalidPathError,
		},
		"path has '/' suffix": {
			path:        "/todo/",
			expect:      "/todo",
			expectError: nil,
		},
		"path doesn't have '/' suffix": {
			path:        "/todo",
			expect:      "/todo",
			expectError: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := normalize(tc.path)
			if tc.expectError != nil {
				assert.Equal(t, tc.expectError, err)
			}
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestBracesIndex(t *testing.T) {
	cases := map[string]struct {
		path        string
		expect      []int
		expectError error
	}{
		"path without braces": {
			path:        "/todo",
			expect:      []int{},
			expectError: nil,
		},
		"path with braces": {
			path:        "/todo/{id}",
			expect:      []int{6, 9},
			expectError: nil,
		},
		"path with multiple param": {
			path:        "/todo/{id}/{field}",
			expect:      []int{6, 9, 11, 17},
			expectError: nil,
		},
		"path with unbalanced braces(braces isn't close)": {
			path:        "/todo/{i{d}",
			expect:      nil,
			expectError: UnbalancedBracesError,
		},
		"path with unbalanced braces(braces starts with '}')": {
			path:        "/todo/}id}",
			expect:      nil,
			expectError: UnbalancedBracesError,
		},
		"path with unbalanced braces('}' is continuous)": {
			path:        "/todo/{id}}",
			expect:      nil,
			expectError: UnbalancedBracesError,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := bracesIndex(tc.path)
			if tc.expectError != nil {
				assert.Equal(t, tc.expectError, err)
			}
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestNewRoute(t *testing.T) {
	cases := map[string]struct {
		path               string
		method             string
		expectParamNames   []string
		expectRegexpString string
	}{
		"path without param": {
			path:               "/todo",
			method:             http.MethodGet,
			expectParamNames:   []string{},
			expectRegexpString: "^/todo[/]?$",
		},
		"path with param": {
			path:               "/todo/{id}",
			method:             http.MethodGet,
			expectParamNames:   []string{"id"},
			expectRegexpString: "^/todo/([^/]+)[/]?$",
		},
		"path with multipule params": {
			path:               "/todo/{id}/{field}",
			method:             http.MethodGet,
			expectParamNames:   []string{"id", "field"},
			expectRegexpString: "^/todo/([^/]+)/([^/]+)[/]?$",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			route, _ := NewRoute(tc.method, tc.path, func(http.ResponseWriter, *http.Request) {})
			assert.Equal(t, tc.method, route.Method)
			assert.Equal(t, tc.path, route.Path)
			assert.Equal(t, tc.expectParamNames, route.ParamNames)
			assert.Equal(t, tc.expectRegexpString, route.Pattern.String())
		})
	}
}

func TestMatch(t *testing.T) {
	cases := map[string]struct {
		method        string
		path          string
		requestMethod string
		requestPath   string
		expectMatch   bool
	}{
		"Route GET /todo match GET /todo": {
			method:        http.MethodGet,
			path:          "/todo",
			requestMethod: http.MethodGet,
			requestPath:   "/todo",
			expectMatch:   true,
		},
		"Route GET /todo match GET /todo/": {
			method:        http.MethodGet,
			path:          "/todo",
			requestMethod: http.MethodGet,
			requestPath:   "/todo/",
			expectMatch:   true,
		},
		"Route GET /todo doesn't match POST /todo": {
			method:        http.MethodGet,
			path:          "/todo",
			requestMethod: http.MethodPost,
			requestPath:   "/todo",
			expectMatch:   false,
		},
		"Route GET /todo/{id} match GET /todo/aaa": {
			method:        http.MethodGet,
			path:          "/todo/{id}",
			requestMethod: http.MethodGet,
			requestPath:   "/todo/aaa",
			expectMatch:   true,
		},
		"Route GET /todo/{id} match GET /todo/123": {
			method:        http.MethodGet,
			path:          "/todo/{id}",
			requestMethod: http.MethodGet,
			requestPath:   "/todo/123",
			expectMatch:   true,
		},
		"Route GET /todo/{id} match GET /todo/123/": {
			method:        http.MethodGet,
			path:          "/todo/{id}",
			requestMethod: http.MethodGet,
			requestPath:   "/todo/123/",
			expectMatch:   true,
		},
		"Route GET /todo/{id}/{field} match GET /todo/aaa/status": {
			method:        http.MethodGet,
			path:          "/todo/{id}/{field}",
			requestMethod: http.MethodGet,
			requestPath:   "/todo/123/status",
			expectMatch:   true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			route, _ := NewRoute(tc.method, tc.path, func(http.ResponseWriter, *http.Request) {})
			req := httptest.NewRequest(tc.requestMethod, fmt.Sprintf("https://example.com%s", tc.requestPath), nil)
			actual := route.Match(req)
			assert.Equal(t, tc.expectMatch, actual)
		})
	}
}
