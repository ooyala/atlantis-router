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
