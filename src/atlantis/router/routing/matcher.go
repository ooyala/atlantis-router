package routing

import (
	"errors"
	"log"
	"net/http"
)

type Matcher interface {
	Match(r *http.Request) bool
}

type matcherMaker func(string) Matcher

type MatcherFactory struct {
	lut map[string]matcherMaker
}

func NewMatcherFactory() *MatcherFactory {
	return &MatcherFactory{
		lut: map[string]matcherMaker{},
	}
}

func (f *MatcherFactory) Register(kind string, maker matcherMaker) {
	if _, ok := f.lut[kind]; ok {
		log.Printf("cannot replace %s in matcher factory", kind)
		return
	}

	f.lut[kind] = maker
}

func (f *MatcherFactory) Make(kind, value string) (Matcher, error) {
	if _, ok := f.lut[kind]; !ok {
		return nil, errors.New("no registered maker")
	}
	return f.lut[kind](value), nil
}

func DefaultMatcherFactory() *MatcherFactory {
	return &MatcherFactory{
		lut: map[string]matcherMaker{
			"static":      NewStaticMatcher,
			"percent":     NewPercentMatcher,
			"host":        NewHostMatcher,
			"multi-host":  NewMultiHostMatcher,
			"header":      NewHeaderMatcher,
			"path-prefix": NewPathPrefixMatcher,
			"path-suffix": NewPathSuffixMatcher,
			"path-regexp": NewPathRegexpMatcher,
		},
	}
}
