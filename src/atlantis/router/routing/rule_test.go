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
