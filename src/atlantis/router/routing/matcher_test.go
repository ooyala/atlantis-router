package routing

import (
	"testing"
)

func TestNewMatcherFactory(t *testing.T) {
	factory := NewMatcherFactory()

	if factory.lut == nil {
		t.Errorf("should allocate lut")
	}
}

func newFalseMatcher(ignore string) Matcher {
	return NewStaticMatcher("false")
}

func newTrueMatcher(ignore string) Matcher {
	return NewStaticMatcher("true")
}

func TestMatcherFactory(t *testing.T) {
	factory := NewMatcherFactory()
	factory.Register("false", newFalseMatcher)
	factory.Register("true", newTrueMatcher)

	falseMatcher, _ := factory.Make("false", "ignored")
	switch falseMatcher.(type) {
	case *StaticMatcher:
		if falseMatcher.(*StaticMatcher).Val != false {
			t.Errorf("should make right kind of matcher")
		}
		break
	default:
		t.Errorf("should make right kind of matcher")
	}

	trueMatcher, _ := factory.Make("true", "ignored")
	switch trueMatcher.(type) {
	case *StaticMatcher:
		if trueMatcher.(*StaticMatcher).Val != true {
			t.Errorf("should make right kind of matcher")
		}
		break
	default:
		t.Errorf("should make right kind of matcher")
	}

	_, err := factory.Make("unregistered", "whaaa!?")
	if err == nil {
		t.Errorf("should return error for unregistered matcher")
	}

	factory.Register("true", newFalseMatcher)
	testMatcher, _ := factory.Make("true", "ignored")
	if testMatcher.(*StaticMatcher).Val != true {
		t.Errorf("should not re-register for kind")
	}
}

func TestDefaultMatcherFactory(t *testing.T) {
	factory := DefaultMatcherFactory()
	falseMatcher, _ := factory.Make("static", "false")
	switch falseMatcher.(type) {
	case *StaticMatcher:
		if falseMatcher.(*StaticMatcher).Val != false {
			t.Errorf("should be pre-populated")
		}
		break
	default:
		t.Errorf("should be pre-populated")
	}
}
