package tinyrouter

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	InvalidPathError      = errors.New("tinyrouter: invalid path")
	UnbalancedBracesError = errors.New("tinyrouter: unbalanced braces")
)

type Route struct {
	ParamNames []string
	Method     string
	Path       string
	Pattern    *regexp.Regexp
	HandleFunc func(http.ResponseWriter, *http.Request)
}

func extractParam(path string, indexList []int) []string {
	paramList := make([]string, 0)
	for i := 0; i < len(indexList); i += 2 {
		paramList = append(paramList, path[indexList[i]+1:indexList[i+1]])
	}
	return paramList
}

func makePatternString(path string, indexList []int) string {
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
	return fmt.Sprintf("^%s[/]?$", builder.String())
}

func makePattern(path string, indexList []int) (*regexp.Regexp, error) {
	pstring := makePatternString(path, indexList)
	return regexp.Compile(pstring)
}

func normalize(raw string) (string, error) {
	if raw[0] != '/' {
		return "", InvalidPathError
	}

	if raw[len(raw)-1] == '/' {
		return raw[:len(raw)-1], nil
	}

	return raw, nil
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

func NewRoute(method string, raw string, f func(http.ResponseWriter, *http.Request)) (*Route, error) {
	path, pathErr := normalize(raw)
	if pathErr != nil {
		return nil, pathErr
	}

	indexList, bracesErr := bracesIndex(path)
	if bracesErr != nil {
		return nil, bracesErr
	}

	pattern, patternErr := makePattern(path, indexList)
	if patternErr != nil {
		return nil, patternErr
	}

	return &Route{
		ParamNames: extractParam(path, indexList),
		Method:     method,
		Path:       raw,
		Pattern:    pattern,
		HandleFunc: f,
	}, nil
}

func (m *Route) Match(r *http.Request) bool {
	matchMethod := m.Method != r.Method
	matchPath := m.Pattern.MatchString(r.URL.Path)

	return matchMethod && matchPath
}
