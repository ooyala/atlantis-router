package routing

import (
	"net/http"
	"testing"
)

func TestHostMatcher(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://white.unicorns.org/magic", nil)

	matcher := NewHostMatcher("white.unicorns.org")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewHostMatcher("pink.unicorns.org")
	if matcher.Match(req) == true {
		t.Errorf("should not match")
	}
}

func TestHeaderMatcher(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://white.unicorns.org/aloha", nil)
	req.Header.Add("unicorn", "rubies")

	matcher := NewHeaderMatcher("unicorn:rubies")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewHeaderMatcher("unicorn:ponies")
	if matcher.Match(req) != false {
		t.Errorf("should not match")
	}
}

func TestHeaderMatcherParseError(t *testing.T) {
	matcher := NewHeaderMatcher("rubies!")

	if matcher.(*HeaderMatcher).ParseError != true {
		t.Errorf("should set parse error")
	}

	matcher = NewHeaderMatcher("rubies!:")

	if matcher.(*HeaderMatcher).ParseError != true {
		t.Errorf("should set parse error")
	}

	matcher = NewHeaderMatcher(":rubies!")

	if matcher.(*HeaderMatcher).ParseError != true {
		t.Errorf("should set parse error")
	}

	if matcher.Match(nil) != false {
		t.Errorf("should not match")
	}
}

func TestAtlantisAppMatcher(t *testing.T) {
	AtlantisAppSuffixes = []string{"unicorns.org", "rainbows.org"}
	quietWhite := NewAtlantisAppMatcher("quiet.white")
	quietWhiteDot := NewAtlantisAppMatcher("quiet.white.")

	req, _ := http.NewRequest("GET", "http://quiet.white.unicorns.org/aloha", nil)
	if quietWhite.Match(req) != true {
		t.Errorf("should match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.white.rainbows.org/aloha", nil)
	if quietWhite.Match(req) != true {
		t.Errorf("should match")
	}

	req, _ = http.NewRequest("GET", "http://quiet.white.rainbowsandunicorns.org/aloha", nil)
	if quietWhite.Match(req) != false {
		t.Errorf("should not match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.white.ugly.unicorns.org/aloha", nil)
	if quietWhite.Match(req) != false {
		t.Errorf("should not match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.ugly.white.unicorns.org/aloha", nil)
	if quietWhite.Match(req) != false {
		t.Errorf("should not match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.white.ugly.rainbows.org/aloha", nil)
	if quietWhite.Match(req) != false {
		t.Errorf("should not match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.ugly.white.rainbows.org/aloha", nil)
	if quietWhite.Match(req) != false {
		t.Errorf("should not match")
	}

	req, _ = http.NewRequest("GET", "http://quiet.white.unicorns.org/aloha", nil)
	if quietWhiteDot.Match(req) != true {
		t.Errorf("should match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.white.rainbows.org/aloha", nil)
	if quietWhiteDot.Match(req) != true {
		t.Errorf("should match")
	}

	req, _ = http.NewRequest("GET", "http://quiet.white.rainbowsandunicorns.org/aloha", nil)
	if quietWhiteDot.Match(req) != false {
		t.Errorf("should not match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.white.ugly.unicorns.org/aloha", nil)
	if quietWhiteDot.Match(req) != false {
		t.Errorf("should not match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.ugly.white.unicorns.org/aloha", nil)
	if quietWhiteDot.Match(req) != false {
		t.Errorf("should not match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.white.ugly.rainbows.org/aloha", nil)
	if quietWhiteDot.Match(req) != false {
		t.Errorf("should not match")
	}
	req, _ = http.NewRequest("GET", "http://quiet.ugly.white.rainbows.org/aloha", nil)
	if quietWhiteDot.Match(req) != false {
		t.Errorf("should not match")
	}
}
