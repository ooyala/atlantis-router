package routing

import (
	"net/http"
	"testing"
)

func TestPathPrefixMatcher(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://in.tet.net/prefix/string?addr=broadcast", nil)

	matcher := NewPathPrefixMatcher("/prefix")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewPathPrefixMatcher("/string")
	if matcher.Match(req) != false {
		t.Errorf("should not match")
	}
}

func TestPathSuffixMatcher(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://in.tet.net/prefix/string?addr=broadcast", nil)

	matcher := NewPathSuffixMatcher("ring")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewPathSuffixMatcher("prefix")
	if matcher.Match(req) != false {
		t.Errorf("should not match")
	}
}

func TestPathRegexpMatcher(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://in.tet.net/prefix/string?addr=broadcast", nil)

	matcher := NewPathRegexpMatcher("pr[e-i]+x/s")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewPathRegexpMatcher("pr[e-ix]+s")
	if matcher.Match(req) != false {
		t.Errorf("should not match")
	}
}

func TestPathRegexpMatcherParseError(t *testing.T) {
	matcher := NewPathRegexpMatcher("pre[e-i")

	if matcher.(*PathRegexpMatcher).ParseError != true {
		t.Errorf("should set parse error")
	}

	if matcher.Match(nil) != false {
		t.Errorf("should not match")
	}
}
