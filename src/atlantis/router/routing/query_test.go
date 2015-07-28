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

func TestQueryParamValueMatcher(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://in.tet.net/path/still_path?param1=value1&param2=value2", nil)

	matcher := NewQueryParamValueMatcher("param1:value1")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewQueryParamValueMatcher("param2:value2")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewQueryParamValueMatcher("param1:value2")
	if matcher.Match(req) != false {
		t.Errorf("should not match")
	}

	matcher = NewQueryParamValueMatcher("param3:value3")
	if matcher.Match(req) != false {
		t.Errorf("should not match")
	}
}

func TestQueryRawRegexpMatcher(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://in.tet.net/path/still_path?param1=value1&param2=value2", nil)

	matcher := NewQueryRawRegexpMatcher("param1=value1($|&)")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewQueryRawRegexpMatcher("[?&]param2=value2($|&)")
	if matcher.Match(req) != true {
		t.Errorf("should match")
	}

	matcher = NewQueryRawRegexpMatcher("[?&]param1=value2($|&)")
	if matcher.Match(req) != false {
		t.Errorf("should not match")
	}

	matcher = NewQueryRawRegexpMatcher("[?&]param3=value3($|&)")
	if matcher.Match(req) != false {
		t.Errorf("should not match")
	}
}
