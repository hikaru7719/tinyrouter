// Copyright (c) 2021 Hikaru Miyahara
// Copyright (c) 2012-2018 The Gorilla Authors.
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

type route struct {
	paramNames []string
	method     string
	path       string
	pattern    *regexp.Regexp
	handleFunc http.HandlerFunc
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

func newRoute(method string, raw string, f func(http.ResponseWriter, *http.Request)) (*route, error) {
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

	return &route{
		paramNames: extractParam(path, indexList),
		method:     method,
		path:       raw,
		pattern:    pattern,
		handleFunc: f,
	}, nil
}

func (m *route) match(r *http.Request) bool {
	matchMethod := m.method == r.Method
	matchPath := m.pattern.MatchString(r.URL.Path)

	return matchMethod && matchPath
}

func (m *route) params(r *http.Request) []string {
	result := m.pattern.FindStringSubmatch(r.URL.Path)
	return result[1:]
}

func (m *route) combineParams(paramValues []string, paramBox map[string]interface{}) {
	if len(m.paramNames) != len(paramValues) {
		return
	}

	for index, name := range m.paramNames {
		paramBox[name] = paramValues[index]
	}
}

func (m *route) setParams(r *http.Request, paramBox map[string]interface{}) {
	if len(m.paramNames) == 0 {
		return
	}
	paramValues := m.params(r)
	m.combineParams(paramValues, paramBox)
}
