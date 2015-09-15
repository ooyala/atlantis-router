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
	"atlantis/router/logger"
	"atlantis/router/routing"
	"net/http"
	"sync"
)

var MaxRoutingHops = 128

type Config struct {
	sync.RWMutex
	MatcherFactory *routing.MatcherFactory
	Pools          map[string]*backend.Pool
	Rules          map[string]*routing.Rule
	Tries          map[string]*routing.Trie
	Ports          map[uint16]*routing.Trie
}

func NewConfig(matcherFactory *routing.MatcherFactory) *Config {
	return &Config{
		MatcherFactory: matcherFactory,
		Pools:          make(map[string]*backend.Pool, 32),
		Rules:          make(map[string]*routing.Rule, 1024),
		Tries:          make(map[string]*routing.Trie, 128),
		Ports:          make(map[uint16]*routing.Trie, 32),
	}
}

// NOTE(manas): this function must be called holding read lock on config
func (c *Config) route(trie *routing.Trie, r *http.Request) *backend.Pool {
	var pool *backend.Pool
	var next *routing.Trie

	next = trie
	for hops := 0; hops < MaxRoutingHops; hops++ {
		pool, next = next.Walk(r)
		if pool != nil {
			return pool
		}
		if next == nil {
			break
		}
	}

	return nil
}

func (c *Config) RouteTrie(trie *routing.Trie, r *http.Request) *backend.Pool {
	c.RLock()
	defer c.RUnlock()

	return c.route(trie, r)
}

func (c *Config) RoutePort(port uint16, r *http.Request) *backend.Pool {
	c.RLock()
	defer c.RUnlock()

	trie, ok := c.Ports[port]
	if ok {
		return c.route(trie, r)
	} else {
		return nil
	}
}

func (c *Config) AddPool(pool Pool) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Pools[pool.Name]; ok {
		logger.Errorf("pool exists in config", pool.Name)
		return
	}

	c.Pools[pool.Name] = c.ConstructPool(pool)

	// update references to this pool
	for _, rule := range c.Rules {
		if rule.Pool == pool.Name {
			rule.PoolPtr = c.Pools[pool.Name]
		}
	}
}

func (c *Config) UpdatePool(pool Pool) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Pools[pool.Name]; !ok {
		logger.Errorf("no pool %s to update", pool.Name)
		return
	}

	c.Pools[pool.Name].Reconfigure(c.ConstructPoolConfig(pool), c.ConstructPoolHeaders(pool))
}

func (c *Config) DelPool(name string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Pools[name]; !ok {
		logger.Errorf("no pool %s to delete", name)
		return
	}

	// nil references to this pool
	for _, rule := range c.Rules {
		if rule.Pool == name {
			rule.PoolPtr = nil
		}
	}

	c.Pools[name].Shutdown()
	delete(c.Pools, name)
}

func (c *Config) AddRule(rule Rule) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Rules[rule.Name]; ok {
		logger.Errorf("rule %s exists in config", rule.Name)
		return
	}

	c.Rules[rule.Name] = c.ConstructRule(rule)

	// update references to this rule
	for _, trie := range c.Tries {
		trie.UpdateRule(c.Rules[rule.Name])
	}
}

func (c *Config) UpdateRule(rule Rule) {
	c.Lock()
	defer c.Unlock()

	c.Rules[rule.Name] = c.ConstructRule(rule)

	// update references to this rule
	for _, trie := range c.Tries {
		trie.UpdateRule(c.Rules[rule.Name])
	}
}

func (c *Config) DelRule(name string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Rules[name]; !ok {
		logger.Errorf("no rule %s to delete", name)
		return
	}

	// nil references to this rule
	dummy := routing.DummyRule(name)
	for _, trie := range c.Tries {
		trie.UpdateRule(dummy)
	}

	delete(c.Rules, name)
}

func (c *Config) AddTrie(trie Trie) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Tries[trie.Name]; ok {
		logger.Errorf("trie %s exists in config", trie.Name)
		return
	}

	c.Tries[trie.Name] = c.ConstructTrie(trie)

	// update references to this trie
	for _, rule := range c.Rules {
		if rule.Next == trie.Name {
			rule.NextPtr = c.Tries[trie.Name]
		}
	}

	for num, _ := range c.Ports {
		if c.Ports[num].Name == trie.Name {
			c.Ports[num] = c.Tries[trie.Name]
		}
	}
}

func (c *Config) UpdateTrie(trie Trie) {
	c.Lock()
	defer c.Unlock()

	c.Tries[trie.Name] = c.ConstructTrie(trie)

	// update references to this trie
	for _, rule := range c.Rules {
		if rule.Next == trie.Name {
			rule.NextPtr = c.Tries[trie.Name]
		}
	}

	for num, _ := range c.Ports {
		if c.Ports[num].Name == trie.Name {
			c.Ports[num] = c.Tries[trie.Name]
		}
	}
}

func (c *Config) DelTrie(name string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Tries[name]; !ok {
		logger.Errorf("no trie %s to delete", name)
		return
	}

	// nil references to this trie
	dummy := routing.DummyTrie(name)
	for _, rule := range c.Rules {
		if rule.Next == name {
			rule.NextPtr = dummy
		}
	}

	for num, _ := range c.Ports {
		if c.Ports[num].Name == name {
			c.Ports[num] = dummy
		}
	}

	delete(c.Tries, name)
}

func (c *Config) AddPort(port Port) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Ports[port.Port]; ok {
		logger.Errorf("port %s exists in config", port.Port)
		return
	}

	trie, ok := c.Tries[port.Trie]
	if !ok {
		trie = routing.DummyTrie(port.Trie)
		logger.Errorf("no trie %s in config", port.Trie)
	}
	c.Ports[port.Port] = trie
}

func (c *Config) UpdatePort(port Port) {
	c.Lock()
	defer c.Unlock()

	trie, ok := c.Tries[port.Trie]
	if !ok {
		trie = routing.DummyTrie(port.Trie)
		logger.Errorf("no trie %s in config", port.Trie)
	}
	c.Ports[port.Port] = trie
}

func (c *Config) DelPort(num uint16) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Ports[num]; !ok {
		logger.Errorf("no port %u to delete", num)
		return
	}

	delete(c.Ports, num)
}
