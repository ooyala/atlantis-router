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
		logger.Errorf("[config %s] %s is not valid duration", name, config.HealthzEvery)
		healthzEvery = defaultHealthzEvery
	}

	healthzTimeout, err := time.ParseDuration(config.HealthzTimeout)
	if err != nil {
		logger.Errorf("[config %s] %s is not valid duration", name, config.HealthzTimeout)
		healthzTimeout = defaultHealthzTimeout
	}

	requestTimeout, err := time.ParseDuration(config.RequestTimeout)
	if err != nil {
		logger.Errorf("[config %s] %s is not valid duration", name, config.RequestTimeout)
		requestTimeout = defaultRequestTimeout
	}

	status := config.Status
	if !backend.IsValidStatus(status) {
		logger.Errorf("[config %s] %s is not valid status", name, config.Status)
		status = "OK"
	}

	return backend.PoolConfig{
		HealthzEvery:   healthzEvery,
		HealthzTimeout: healthzTimeout,
		RequestTimeout: requestTimeout,
		Status:         status,
	}
}

func (c *Config) ConstructPoolHeaders(pool Pool) (httpHeaders map[string]string) {
	httpHeaders = make(map[string]string)
	for _, hdr := range pool.Headers {
		if &hdr != nil {
			httpHeaders[hdr.Key] = hdr.Value
		}
	}
	return httpHeaders
}

func (c *Config) ConstructPool(pool Pool) *backend.Pool {
	return backend.NewPool(pool.Name, c.ConstructPoolConfig(pool), c.ConstructPoolHeaders(pool))
}

func (c *Config) ConstructRule(rule Rule) *routing.Rule {
	if rule.Next == "" && rule.Pool == "" {
		logger.Errorf("[rule %s] no pool or trie", rule.Name)
		return routing.DummyRule(rule.Name)
	}

	var next *routing.Trie
	if rule.Next != "" {
		next = c.Tries[rule.Next]
		if next == nil {
			logger.Errorf("[rule %s] trie %s absent", rule.Name, rule.Next)
			next = routing.DummyTrie(rule.Next)
		}
	}

	var pool *backend.Pool
	if rule.Pool != "" {
		pool = c.Pools[rule.Pool]
		if pool == nil {
			logger.Errorf("[rule %s] trie %s absent", rule.Name, rule.Pool)
			pool = backend.DummyPool(rule.Pool)
		}
	}

	matcher, err := c.MatcherFactory.Make(rule.Type, rule.Value)
	if err != nil {
		logger.Errorf("[rule %s] setting matcher false", rule.Name)
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
			logger.Errorf("[trie %s] rule %s absent", trie.Name, rule)
			list = append(list, routing.DummyRule(rule))
		}
	}

	return routing.NewTrie(trie.Name, list)
}
