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
