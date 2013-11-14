package routing

import (
	"net/http"
	"regexp"
	"strings"
)

type PathPrefixMatcher struct {
	Prefix string
}

func (p *PathPrefixMatcher) Match(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, p.Prefix) {
		return true
	}
	return false
}

func NewPathPrefixMatcher(prefix string) Matcher {
	return &PathPrefixMatcher{prefix}
}

type PathSuffixMatcher struct {
	Suffix string
}

func (p *PathSuffixMatcher) Match(r *http.Request) bool {
	if strings.HasSuffix(r.URL.Path, p.Suffix) {
		return true
	}
	return false
}

func NewPathSuffixMatcher(suffix string) Matcher {
	return &PathSuffixMatcher{suffix}
}

type PathRegexpMatcher struct {
	Regexp     *regexp.Regexp
	ParseError bool
}

func (p *PathRegexpMatcher) Match(r *http.Request) bool {
	if p.ParseError {
		return false
	}

	return p.Regexp.MatchString(r.URL.Path)
}

func NewPathRegexpMatcher(data string) Matcher {
	regexp, err := regexp.Compile(data)
	if err != nil {
		return &PathRegexpMatcher{nil, true}
	}

	return &PathRegexpMatcher{regexp, false}
}
