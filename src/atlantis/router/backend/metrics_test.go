package backend

import (
	"testing"
)

func TestNewMetrics(t *testing.T) {
	metrics := NewServerMetrics()
	if metrics.RequestsInFlight != 0 {
		t.Errorf("should start with 0 requests in flight")
	}
	if metrics.RequestsServiced != 0 {
		t.Errorf("should start with 0 requests serviced")
	}
	return
}

func TestRequestStart(t *testing.T) {
	metrics := NewServerMetrics()
	metrics.RequestStart()
	if metrics.RequestsInFlight != 1 {
		t.Errorf("should increment requests in flight")
	}
	if metrics.RequestsServiced != 1 {
		t.Errorf("should increment requests serviced")
	}
	return
}

func TestRequestDone(t *testing.T) {
	metrics := NewServerMetrics()
	metrics.RequestStart()
	metrics.RequestDone()
	if metrics.RequestsInFlight != 0 {
		t.Errorf("should decrement requests in flight")
	}
	if metrics.RequestsServiced != 1 {
		t.Errorf("should not decrement requests serviced")
	}
}

func TestCost(t *testing.T) {
	N := 10

	metrics := NewServerMetrics()
	for i := 0; i < N; i++ {
		metrics.RequestStart()
	}
	if metrics.Cost() != uint32(N) {
		t.Errorf("should report requests in flight")
	}

	for i := 0; i < N; i++ {
		metrics.RequestDone()
	}
	if metrics.Cost() != 0 {
		t.Errorf("should report requests in flight")
	}
}
