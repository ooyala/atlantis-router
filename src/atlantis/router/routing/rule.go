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
