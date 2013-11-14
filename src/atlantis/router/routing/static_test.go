package routing

import (
	"testing"
)

func TestStaticMatcher(t *testing.T) {
	matcherT := NewStaticMatcher("true")
	if matcherT.Match(nil) != true {
		t.Errorf("should match true")
	}

	matcherF := NewStaticMatcher("false")
	if matcherF.Match(nil) != false {
		t.Errorf("should match false")
	}
}

func TestStaticMatcherParseError(t *testing.T) {
	matcher := NewStaticMatcher("maybe")

	if matcher.(*StaticMatcher).ParseError != true {
		t.Errorf("should set parse error")
	}

	if matcher.Match(nil) != false {
		t.Errorf("should not match")
	}
}

func TestPercentMatcher(t *testing.T) {
	matcher := NewPercentMatcher("10.0")

	mcount := 0
	for i := 0; i < 1000; i++ {
		if matcher.Match(nil) {
			mcount++
		}
	}

	// statistics, thou art a heartless bitch
	if mcount > 131 || mcount < 69 {
		t.Errorf("should match 1 in 10 requests")
	}
}

func TestPercentMatcherParseError(t *testing.T) {
	matcher := NewPercentMatcher("español")

	if matcher.(*PercentMatcher).ParseError != true {
		t.Errorf("should set parse error")
	}

	matcher = NewPercentMatcher("-10.0")

	if matcher.(*PercentMatcher).ParseError != true {
		t.Errorf("should set parse error")
	}

	if matcher.Match(nil) != false {
		t.Errorf("should not match")
	}
}
