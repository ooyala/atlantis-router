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
	"errors"
	"log"
	"net/http"
)

type Matcher interface {
	Match(r *http.Request) bool
}

type matcherMaker func(string) Matcher

type MatcherFactory struct {
	lut map[string]matcherMaker
}

func NewMatcherFactory() *MatcherFactory {
	return &MatcherFactory{
		lut: map[string]matcherMaker{},
	}
}

func (f *MatcherFactory) Register(kind string, maker matcherMaker) {
	if _, ok := f.lut[kind]; ok {
		log.Printf("cannot replace %s in matcher factory", kind)
		return
	}

	f.lut[kind] = maker
}

func (f *MatcherFactory) Make(kind, value string) (Matcher, error) {
	if _, ok := f.lut[kind]; !ok {
		return nil, errors.New("no registered maker")
	}
	return f.lut[kind](value), nil
}

func DefaultMatcherFactory() *MatcherFactory {
	return &MatcherFactory{
		lut: map[string]matcherMaker{
			"static":      NewStaticMatcher,
			"percent":     NewPercentMatcher,
			"host":        NewHostMatcher,
			"multi-host":  NewMultiHostMatcher,
			"header":      NewHeaderMatcher,
			"path-prefix": NewPathPrefixMatcher,
			"path-suffix": NewPathSuffixMatcher,
			"path-regexp": NewPathRegexpMatcher,
			"query-param-value": NewQueryParamValueMatcher,
			"query-raw-regexp": NewQueryRawRegexpMatcher,
		},
	}
}
