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
	"strings"
)

type HostMatcher struct {
	Host string
}

func (h *HostMatcher) Match(r *http.Request) bool {
	return strings.SplitN(r.Host, ":", 1)[0] == h.Host
}

func NewHostMatcher(r string) Matcher {
	return &HostMatcher{r}
}

type MultiHostMatcher struct {
	Hosts []string
}

func (m *MultiHostMatcher) Match(r *http.Request) bool {
	for _, root := range m.Hosts {
		if r.Host == root {
			return true
		}
	}
	return false
}

func NewMultiHostMatcher(r string) Matcher {
	m := &MultiHostMatcher{
		Hosts: []string{},
	}
	name := strings.Split(r, ":")[0]
	doms := strings.Split(r, ":")[1]
	for _, dom := range strings.Split(doms, ",") {
		m.Hosts = append(m.Hosts, name+"."+dom)
	}
	return m
}

type HeaderMatcher struct {
	Header     string
	Value      string
	ParseError bool
}

func (h *HeaderMatcher) Match(r *http.Request) bool {
	if h.ParseError {
		return false
	}
	return r.Header.Get(h.Header) == h.Value
}

func NewHeaderMatcher(r string) Matcher {
	hdrVal := strings.Split(r, ":")
	if len(hdrVal) != 2 {
		return &HeaderMatcher{"", "", true}
	}

	if hdrVal[0] == "" || hdrVal[1] == "" {
		return &HeaderMatcher{"", "", true}
	}

	return &HeaderMatcher{hdrVal[0], hdrVal[1], false}
}
