package config

import (
	"atlantis/router/backend"
	"atlantis/router/logger"
	"atlantis/router/routing"
	"net/http"
	"sync"
)

type Config struct {
	sync.RWMutex
	MatcherFactory *routing.MatcherFactory
	Pools          map[string]*backend.Pool
	Rules          map[string]*routing.Rule
	Tries          map[string]*routing.Trie
}

func NewConfig(matcherFactory *routing.MatcherFactory) *Config {
	return &Config{
		MatcherFactory: matcherFactory,
		Pools:          make(map[string]*backend.Pool, 20),
		Rules:          make(map[string]*routing.Rule, 1000),
		Tries:          make(map[string]*routing.Trie, 100),
	}
}

func (c *Config) RouteFrom(trie string, r *http.Request) *backend.Pool {
	c.RLock()
	defer c.RUnlock()

	var pool *backend.Pool

	next := c.Tries[trie]
	for next != nil {
		pool, next = next.Walk(r)
		if pool != nil {
			return pool
		}
	}
	return nil
}

func (c *Config) Route(r *http.Request) *backend.Pool {
	return c.RouteFrom("root", r)
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

	c.Pools[pool.Name].Reconfigure(c.ConstructPoolConfig(pool))
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

	if _, ok := c.Rules[rule.Name]; ok {
		delete(c.Rules, rule.Name)
	}

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
}

func (c *Config) UpdateTrie(trie Trie) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Tries[trie.Name]; ok {
		delete(c.Tries, trie.Name)
	}

	c.Tries[trie.Name] = c.ConstructTrie(trie)

	// update references to this trie
	for _, rule := range c.Rules {
		if rule.Next == trie.Name {
			rule.NextPtr = c.Tries[trie.Name]
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

	delete(c.Tries, name)
}
