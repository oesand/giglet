package url

import (
	"regexp"
)

var (
	paramRegex = regexp.MustCompile("{([a-zA-Z_][a-zA-Z0-9_]*)}")
	defaultPatternRegex = regexp.MustCompile("^/?$")
)

func compileUrlPattern(pattern string) *regexp.Regexp {
	if len(pattern) == 0 {
		panic("route pattern cannot be empty")
	} else if pattern[0] != '/' {
		panic("pattern \""+pattern+"\" should starts with \"/\"")
	}

	last := pattern[len(pattern) - 1]
	if len(pattern) > 1 && last == '/' {
		panic("pattern \""+pattern+"\" cannot ends with \"/\"")
	}

	if pattern == "/" {
		return defaultPatternRegex
	}

	parsed := "^" + paramRegex.ReplaceAllStringFunc(pattern, func(p string) string {
		name := p[1:len(p) - 1]
		return `(?P<`+name+`>[^/]+)`;
	})
	if last != '*' {
		parsed += "/?$"
	}
	return regexp.MustCompile(parsed)
}

type Route[ST any] struct {
	path string
	pattern *regexp.Regexp
	names []string

	stored *ST
}

func createRoute[ST any](pattern string, stored *ST) *Route[ST] {
	match := compileUrlPattern(pattern)
	return &Route[ST]{
		path: pattern,
		pattern: match,
		names: match.SubexpNames()[1:],
		stored: stored,
	}
}

func (current *Route[ST]) IsMatch(path string) bool {
	return current.pattern.MatchString(path)
}

func (current *Route[ST]) Names() []string {
	return current.names
}

func (current *Route[ST]) ParseValues(path string) []string {
	return current.pattern.FindStringSubmatch(path)[1:]
}

func (current *Route[ST]) VisitParams(path string, handler func (key, value string)) {
	values := current.ParseValues(path)
	for i, name := range current.names {
		handler(name, values[i])
	}
}

func (current *Route[ST]) GenerateParams(path string) map[string]string {
	values := current.ParseValues(path)
	store := make(map[string]string, len(values))
	for i, name := range current.names {
		store[name] = values[i]
	}
	return store
}

func (current *Route[ST]) Stored() *ST {
	return current.stored
}
