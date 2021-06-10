package tinyrouter

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	UnbalancedBracesError = errors.New("tinyrouter: unbalanced braces")
)

type Route struct {
	ParamName  []string
	Method     string
	Path       string
	Pattern    regexp.Regexp
	HandleFunc func(http.ResponseWriter, *http.Request)
}

func extractParam(path string, indexList []int) []string {
	paramList := make([]string, 0)
	for i := 0; i < len(indexList); i += 2 {
		paramList = append(paramList, path[indexList[i]+1:indexList[i+1]])
	}
	return paramList
}

func makePattern(path string, indexList []int) string {
	var builder strings.Builder
	if len(indexList) == 0 {
		builder.WriteString(path)
	} else {
		var start int = 0
		var next int
		for i := 0; i < len(indexList); i += 2 {
			next = indexList[i]
			builder.WriteString(path[start:next])
			builder.WriteString("([^/]+)")
			start = indexList[i+1] + 1
		}
	}
	return fmt.Sprintf("^%s$", builder.String())
}

func bracesIndex(path string) ([]int, error) {
	indexList := make([]int, 0)
	checkBraces := 0
	for i := 0; i < len(path); i++ {
		switch path[i] {
		case '{':
			if checkBraces++; checkBraces != 1 {
				return nil, UnbalancedBracesError
			}
			indexList = append(indexList, i)
		case '}':
			if checkBraces--; checkBraces != 0 {
				return nil, UnbalancedBracesError
			}
			indexList = append(indexList, i)
		}
	}
	return indexList, nil
}

func NewRoute(method string, path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return &Route{
		Method:     method,
		Path:       path,
		HandleFunc: f,
	}
}

func (m *Route) Match(r *http.Request) bool {
	return false
}
