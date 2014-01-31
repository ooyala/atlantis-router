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
)

type Rule struct {
	Name    string
	Dummy   bool
	Matcher Matcher
	Next    string
	NextPtr *Trie
	Pool    string
	PoolPtr *backend.Pool
}

func DummyRule(name string) *Rule {
	return &Rule{
		Name:  name,
		Dummy: true,
	}
}

func NewRule(name string, matcher Matcher, next *Trie, pool *backend.Pool) *Rule {
	var nextName string
	if next != nil {
		nextName = next.Name
	}

	var poolName string
	if pool != nil {
		poolName = pool.Name
	}

	return &Rule{
		Name:    name,
		Dummy:   false,
		Matcher: matcher,
		Next:    nextName,
		NextPtr: next,
		Pool:    poolName,
		PoolPtr: pool,
	}
}
