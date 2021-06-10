package tinyrouter

import (
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
