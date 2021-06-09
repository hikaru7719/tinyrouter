package tinyrouter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakePattern(t *testing.T) {
	cases := map[string]struct {
		path   string
		expect string
	}{
		"path without param": {
			path:   "/todo",
			expect: "^/todo$",
		},
		"path with param": {
			path:   "/todo/{id}",
			expect: "^/todo/([^/]+)$",
		},
		"path with multiple param": {
			path:   "/todo/{id}/{field}",
			expect: "^/todo/([^/]+)/([^/]+)$",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := makePattern(tc.path)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
