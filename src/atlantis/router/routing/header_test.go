/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

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

	req, _ = http.NewRequest("GET", "http://white.unicorns.org/magic:8080", nil)

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

func TestMultiHostMatcher(t *testing.T) {
	quietWhite := NewMultiHostMatcher("quiet.white:unicorns.org,rainbows.org")

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
}
