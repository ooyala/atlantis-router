package routing

import (
	"math/rand"
	"net/http"
	"strconv"
)

type StaticMatcher struct {
	Val        bool
	ParseError bool
}

func (s *StaticMatcher) Match(r *http.Request) bool {
	if s.ParseError {
		return false
	}
	return s.Val
}

func NewStaticMatcher(data string) Matcher {
	val, err := strconv.ParseBool(data)
	if err != nil {
		return &StaticMatcher{false, true}
	}

	return &StaticMatcher{val, false}
}

type PercentMatcher struct {
	Fraction   float64
	ParseError bool
}

func (p *PercentMatcher) Match(r *http.Request) bool {
	if p.ParseError {
		return false
	}
	return rand.Float64() < p.Fraction
}

func NewPercentMatcher(data string) Matcher {
	val, err := strconv.ParseFloat(data, 64)
	if err != nil || val < 0.0 {
		return &PercentMatcher{0, true}
	}

	return &PercentMatcher{val / 100.0, false}
}
