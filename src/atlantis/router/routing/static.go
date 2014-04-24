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
	"math/rand"
	"net/http"
	"strconv"
)

type StaticMatcher struct {
	Val        bool
	ParseError bool
}

func (s *StaticMatcher) Match(r *http.Request) bool {
	if s.ParseError {
		return false
	}
	return s.Val
}

func NewStaticMatcher(data string) Matcher {
	val, err := strconv.ParseBool(data)
	if err != nil {
		return &StaticMatcher{false, true}
	}

	return &StaticMatcher{val, false}
}

type PercentMatcher struct {
	Fraction   float64
	ParseError bool
}

func (p *PercentMatcher) Match(r *http.Request) bool {
	if p.ParseError {
		return false
	}
	return rand.Float64() < p.Fraction
}

func NewPercentMatcher(data string) Matcher {
	val, err := strconv.ParseFloat(data, 64)
	if err != nil || val < 0.0 {
		return &PercentMatcher{0, true}
	}

	return &PercentMatcher{val / 100.0, false}
}
