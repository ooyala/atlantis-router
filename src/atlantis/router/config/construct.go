package config

import (
	"atlantis/router/backend"
	"atlantis/router/routing"
	"log"
	"time"
)

const (
	defaultHealthzEvery   = 1 * time.Minute
	defaultHealthzTimeout = 9 * time.Second
	defaultRequestTimeout = 1 * time.Minute
)

func (c *Config) ConstructServer(host Host) *backend.Server {
	return backend.NewServer(host.Address)
}

func (c *Config) ConstructPoolConfig(config PoolConfig) backend.PoolConfig {
	healthzEvery, err := time.ParseDuration(config.HealthzEvery)
	if err != nil {
		log.Printf("cannot parse %s as duration", config.HealthzEvery)
		healthzEvery = defaultHealthzEvery
	}

	healthzTimeout, err := time.ParseDuration(config.HealthzTimeout)
	if err != nil {
		log.Printf("cannot parse %s as duration", config.HealthzTimeout)
		healthzTimeout = defaultHealthzTimeout
	}

	requestTimeout, err := time.ParseDuration(config.RequestTimeout)
	if err != nil {
		log.Printf("cannot parse %s as duration", config.RequestTimeout)
		requestTimeout = defaultRequestTimeout
	}

	status := config.Status
	if !backend.IsValidStatus(status) {
		log.Printf("%s is not valid status", config.Status)
		status = "OK"
	}

	return backend.PoolConfig{
		HealthzEvery:   healthzEvery,
		HealthzTimeout: healthzTimeout,
		RequestTimeout: requestTimeout,
		Status:         status,
	}
}

func (c *Config) ConstructPool(pool Pool) *backend.Pool {
	return backend.NewPool(pool.Name, c.ConstructPoolConfig(pool.Config))
}

func (c *Config) ConstructRule(rule Rule) *routing.Rule {
	if rule.Next == "" && rule.Pool == "" {
		log.Printf("no pool or trie in rule %s", rule.Name)
		return routing.DummyRule(rule.Name)
	}

	var next *routing.Trie
	if rule.Next != "" {
		next = c.Tries[rule.Next]
		if next == nil {
			log.Printf("no trie %s for rule %s", rule.Next, rule.Name)
			next = routing.DummyTrie(rule.Next)
		}
	}

	var pool *backend.Pool
	if rule.Pool != "" {
		pool = c.Pools[rule.Pool]
		if pool == nil {
			log.Printf("no pool %s for rule %s", rule.Pool, rule.Name)
			pool = backend.DummyPool(rule.Pool)
		}
	}

	matcher, err := c.MatcherFactory.Make(rule.Type, rule.Value)
	if err != nil {
		log.Printf("error, replacing %s with false", rule.Type)
		matcher = routing.NewStaticMatcher("false")
	}

	return routing.NewRule(rule.Name, matcher, next, pool)
}

func (c *Config) ConstructTrie(trie Trie) *routing.Trie {
	list := []*routing.Rule{}

	for _, rule := range trie.Rules {
		if _, ok := c.Rules[rule]; ok {
			list = append(list, c.Rules[rule])
		} else {
			log.Printf("no rule %s for trie %s", rule, trie.Name)
			list = append(list, routing.DummyRule(rule))
		}
	}

	return routing.NewTrie(trie.Name, list)
}
