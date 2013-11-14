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
