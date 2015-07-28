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
	"net/url"
	"regexp"
	"strings"
)

type QueryParamValueMatcher struct {
	Param      string
	Value      string
	ParseError bool
}

func (p *QueryParamValueMatcher) Match(r *http.Request) bool {
	if p.ParseError {
		return false
	}

	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return false
	}

	pvals, ok := values[p.Param]
	if !ok {
		return false
	}

	for _, pval := range pvals {
		if pval == p.Value {
			return true
		}
	}

	return false
}

func NewQueryParamValueMatcher(r string) Matcher {
	paramVal := strings.Split(r, ":")
	if len(paramVal) != 2 {
		return &QueryParamValueMatcher{"", "", true}
	}

	if paramVal[0] == "" || paramVal[1] == "" {
		return &QueryParamValueMatcher{"", "", true}
	}

	return &QueryParamValueMatcher{paramVal[0], paramVal[1], false}
}

type QueryRawRegexpMatcher struct {
	Regexp     *regexp.Regexp
	ParseError bool
}

func (p *QueryRawRegexpMatcher) Match(r *http.Request) bool {
	if p.ParseError {
		return false
	}

	return p.Regexp.MatchString(r.URL.RawQuery)
}

func NewQueryRawRegexpMatcher(data string) Matcher {
	regexp, err := regexp.Compile(data)
	if err != nil {
		return &QueryRawRegexpMatcher{nil, true}
	}

	return &QueryRawRegexpMatcher{regexp, false}
}
