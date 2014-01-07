package config

import (
	"atlantis/router/backend"
	"atlantis/router/logger"
	"atlantis/router/routing"
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

func (c *Config) ConstructPoolConfig(pool Pool) backend.PoolConfig {
	name, config := pool.Name, pool.Config

	healthzEvery, err := time.ParseDuration(config.HealthzEvery)
	if err != nil {
		logger.Errorf("[CONFIG %s] %s is not valid duration", name, config.HealthzEvery)
		healthzEvery = defaultHealthzEvery
	}

	healthzTimeout, err := time.ParseDuration(config.HealthzTimeout)
	if err != nil {
		logger.Errorf("[CONFIG %s] %s is not valid duration", name, config.HealthzTimeout)
		healthzTimeout = defaultHealthzTimeout
	}

	requestTimeout, err := time.ParseDuration(config.RequestTimeout)
	if err != nil {
		logger.Errorf("[CONFIG %s] %s is not valid duration", name, config.RequestTimeout)
		requestTimeout = defaultRequestTimeout
	}

	status := config.Status
	if !backend.IsValidStatus(status) {
		logger.Errorf("[CONFIG %s] %s is not valid status", name, config.Status)
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
	return backend.NewPool(pool.Name, c.ConstructPoolConfig(pool))
}

func (c *Config) ConstructRule(rule Rule) *routing.Rule {
	if rule.Next == "" && rule.Pool == "" {
		logger.Errorf("[RULE %s] no pool or trie", rule.Name)
		return routing.DummyRule(rule.Name)
	}

	var next *routing.Trie
	if rule.Next != "" {
		next = c.Tries[rule.Next]
		if next == nil {
			logger.Errorf("[RULE %s] trie %s absent", rule.Name, rule.Next)
			next = routing.DummyTrie(rule.Next)
		}
	}

	var pool *backend.Pool
	if rule.Pool != "" {
		pool = c.Pools[rule.Pool]
		if pool == nil {
			logger.Errorf("[RULE %s] trie %s absent", rule.Name, rule.Pool)
			pool = backend.DummyPool(rule.Pool)
		}
	}

	matcher, err := c.MatcherFactory.Make(rule.Type, rule.Value)
	if err != nil {
		logger.Errorf("[RULE %s] setting matcher false", rule.Name)
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
			logger.Errorf("[TRIE %s] rule %s absent", trie.Name, rule)
			list = append(list, routing.DummyRule(rule))
		}
	}

	return routing.NewTrie(trie.Name, list)
}
