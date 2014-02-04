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
	"atlantis/router/backend"
	"net/http"
)

type Trie struct {
	Name  string
	Dummy bool
	List  []*Rule
}

func DummyTrie(name string) *Trie {
	return &Trie{
		Name:  name,
		Dummy: true,
		List:  []*Rule{},
	}
}

func NewTrie(name string, list []*Rule) *Trie {
	return &Trie{
		Name:  name,
		Dummy: false,
		List:  list,
	}
}

func (t *Trie) UpdateRule(update *Rule) {
	for i, rule := range t.List {
		if rule.Name == update.Name {
			t.List[i] = update
		}
	}
}

func (t *Trie) Walk(r *http.Request) (*backend.Pool, *Trie) {
	for _, rule := range t.List {
		if rule.Dummy {
			continue
		}
		if rule.Matcher.Match(r) {
			return rule.PoolPtr, rule.NextPtr
		}
	}
	return nil, nil
}
