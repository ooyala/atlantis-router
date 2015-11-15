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

package config

import (
	"atlantis/router/backend"
	"atlantis/router/routing"
	"testing"
)

func TestConstructServer(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	server := config.ConstructServer(Host{
		Address: "localhost:8080",
	})

	if server.Address != "localhost:8080" {
		t.Errorf("should construct server accurately")
	}
}

func TestConstructPoolConfig(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	test := Pool{
		Name: "Aloha!",
		Config: PoolConfig{
			HealthzEvery:   "Saturn",
			HealthzTimeout: "Jupiter",
			RequestTimeout: "Mars",
			Status:         "Excellent",
		},
	}

	parsed := config.ConstructPoolConfig(test)

	if parsed.HealthzEvery == 0 || parsed.HealthzTimeout == 0 || parsed.RequestTimeout == 0 ||
		!backend.IsValidStatus(parsed.Status) {
		t.Errorf("should default to sane defaults")
	}
}

func TestConstructPoolHeaders(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())
	headers := make([]HttpHeader, 1)
	headers[0] = HttpHeader{Key: "Cache-Control", Value: "max-age:3600"}
	test := Pool{
		Name: "Aloha!",
		Config: PoolConfig{
			HealthzEvery:   "Saturn",
			HealthzTimeout: "Jupiter",
			RequestTimeout: "Mars",
			Status:         "Excellent",
		},
		Headers: headers,
	}

	parsed := config.ConstructPoolHeaders(test)
	val, ok := parsed["Cache-Control"]
	if !ok {
		t.Errorf("Failed to contruct Headers")
	}
	if val != "max-age:3600" {
		t.Errorf("Failed to contruct Headers. Invalid value expected max-age:3600 and got ", val)
	}
}

func TestEmptyConstructPoolHeaders(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())
	test := Pool{
		Name: "Aloha!",
		Config: PoolConfig{
			HealthzEvery:   "Saturn",
			HealthzTimeout: "Jupiter",
			RequestTimeout: "Mars",
			Status:         "Excellent",
		},
	}

	parsed := config.ConstructPoolHeaders(test)
	if parsed == nil {
		t.Errorf("nil Headers")
	}
}

func TestConstructRuleEmpty(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	test := Rule{
		Name:  "test",
		Type:  "host",
		Value: "www.ooyala.com",
		Next:  "",
		Pool:  "",
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore empty rule")
		}
	}()
	parsed := config.ConstructRule(test)

	if parsed.Dummy != true {
		t.Errorf("should return dummy rule")
	}
}

func TestConstructRuleBadPool(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())
	config.AddTrie(meatTrie())

	test := Rule{
		Name:  "test",
		Type:  "host",
		Value: "www.ooyala.com",
		Next:  "meatTrie",
		Pool:  "butcheryPool",
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent pools")
		}
	}()
	parsed := config.ConstructRule(test)

	if parsed.PoolPtr.Dummy != true {
		t.Errorf("should use dummy rule")
	}
}

func TestConstructRuleBadNext(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())
	config.AddPool(butcheryPool())

	test := Rule{
		Name:  "test",
		Type:  "host",
		Value: "www.ooyala.com",
		Next:  "meatTrie",
		Pool:  "butcherPool",
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent next")
		}
	}()
	parsed := config.ConstructRule(test)

	if parsed.NextPtr.Dummy != true {
		t.Errorf("should use dummy trie")
	}
}

func TestConstructRuleBadMatcher(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())
	config.AddPool(butcheryPool())
	config.AddTrie(meatTrie())

	test := Rule{
		Name:  "test",
		Type:  "anyhow",
		Value: "www.ooyala.com",
		Next:  "meatTrie",
		Pool:  "butcheryPool",
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore bad rules")
		}
	}()
	config.ConstructRule(test)
}

func TestConstructRule(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())
	config.AddPool(butcheryPool())
	config.AddTrie(meatTrie())

	test := Rule{
		Name:  "test",
		Type:  "anyhow",
		Value: "www.ooyala.com",
		Next:  "meatTrie",
		Pool:  "butcheryPool",
	}

	parsed := config.ConstructRule(test)

	if parsed.NextPtr != config.Tries["meatTrie"] || parsed.PoolPtr != config.Pools["butcheryPool"] {
		t.Errorf("should construct rule accurately")
	}
}
