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
