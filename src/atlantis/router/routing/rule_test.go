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
	"testing"
)

func TestDummyRule(t *testing.T) {
	rule := DummyRule("test")

	if rule.Name != "test" {
		t.Errorf("should set name")
	}
	if rule.Dummy != true {
		t.Errorf("should set dummy")
	}
}

func TestNewRule(t *testing.T) {
	trie := DummyTrie("test")
	if trie == nil {
		t.Fatalf("cannot create dummy trie")
	}

	pool := backend.DummyPool("test")
	if pool == nil {
		t.Fatalf("cannot create dummy pool")
	}

	rule := NewRule("test", NewStaticMatcher("false"), trie, pool)
	if rule.Dummy != false {
		t.Errorf("should not set dummy")
	}
	if rule.Next != "test" {
		t.Errorf("should set next")
	}
	if rule.NextPtr != trie {
		t.Errorf("should set next ptr")
	}
	if rule.Pool != "test" {
		t.Errorf("should set pool")
	}
	if rule.PoolPtr != pool {
		t.Errorf("should set pool ptr")
	}
}
