package routing

import (
	"atlantis/router/backend"
	"net/http"
	"testing"
)

func TestDummyTrie(t *testing.T) {
	trie := DummyTrie("test")

	if trie.Name != "test" {
		t.Errorf("should set name")
	}
	if trie.Dummy != true {
		t.Errorf("should set dummy")
	}
}

func TestNewTrie(t *testing.T) {
	rule0 := DummyRule("rule0")
	rule1 := DummyRule("rule1")
	if rule0 == nil || rule1 == nil {
		t.Fatalf("cannot create dummy rule")
	}
	rules := []*Rule{rule0, rule1}

	trie := NewTrie("test", rules)
	if trie.Name != "test" {
		t.Errorf("should set name")
	}
	if trie.List[0] != rules[0] || trie.List[1] != rules[1] {
		t.Errorf("should set list")
	}
}

func TestUpdateTrie(t *testing.T) {
	rule0 := DummyRule("rule0")
	rule1 := DummyRule("rule1")
	rules := []*Rule{rule0, rule1}

	trie := NewTrie("test", rules)
	rule := DummyRule("rule0")

	trie.UpdateRule(rule)
	if trie.List[0] != rule {
		t.Errorf("should update rule")
	}
}

func TestWalk(t *testing.T) {
	matchT := NewStaticMatcher("true")
	matchF := NewStaticMatcher("false")

	trie0 := DummyTrie("test")
	trie1 := DummyTrie("test")
	trie2 := DummyTrie("test")
	if trie0 == nil || trie1 == nil || trie2 == nil {
		t.Errorf("cannot create dummy trie")
	}

	pool0 := backend.DummyPool("test")
	pool1 := backend.DummyPool("test")
	pool2 := backend.DummyPool("test")
	if pool0 == nil || pool1 == nil || pool2 == nil {
		t.Errorf("cannot create dummy pool")
	}

	rule0 := NewRule("rule0", matchF, trie0, pool0)
	rule1 := NewRule("rule1", matchT, trie1, pool1)
	rule2 := NewRule("rule2", matchT, trie2, pool2)

	rules := []*Rule{rule0, rule1, rule2}
	trie := NewTrie("test", rules)

	req, _ := http.NewRequest("GET", "/", nil)
	pool, next := trie.Walk(req)
	if pool != pool1 || next != trie1 {
		t.Errorf("should return first match")
	}
}

func TestWalkDummy(t *testing.T) {
	matchT := NewStaticMatcher("true")
	matchF := NewStaticMatcher("false")

	trie0 := DummyTrie("test")
	trie1 := DummyTrie("test")
	if trie0 == nil || trie1 == nil {
		t.Errorf("cannot create dummy trie")
	}

	pool0 := backend.DummyPool("test")
	pool1 := backend.DummyPool("test")
	if pool0 == nil || pool1 == nil {
		t.Errorf("cannot create dummy pool")
	}

	rule0 := NewRule("rule0", matchT, trie0, pool0)
	rule0.Dummy = true

	rule1 := NewRule("rule1", matchF, trie1, pool1)

	rules := []*Rule{rule0, rule1}
	trie := NewTrie("test", rules)

	req, _ := http.NewRequest("GET", "/", nil)
	pool, next := trie.Walk(req)
	if pool != nil || next != nil {
		t.Errorf("should not match dummy")
	}
}
