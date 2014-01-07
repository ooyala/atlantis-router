package backend

import (
	"atlantis/router/testutils"
	"net/http"
	"testing"
	"time"
)

func TestNewServerStatus(t *testing.T) {
	tstart := time.Now()

	status := NewServerStatus()
	if status.Current != StatusMaintenance {
		t.Errorf("should set status to maintenance")
	}

	if status.checked.UnixNano() < tstart.UnixNano() ||
		status.changed.UnixNano() < tstart.UnixNano() {
		t.Errorf("should set checked and changed")
	}
}

func TestSet(t *testing.T) {
	status := NewServerStatus()
	tcreate := status.changed

	status.Set(StatusCritical)
	if status.checked.UnixNano() <= tcreate.UnixNano() ||
		status.changed.UnixNano() <= tcreate.UnixNano() {
		t.Errorf("should set checked and changed when changes")
	}
	tmodify := status.checked

	status.Set(StatusCritical)
	if status.checked.UnixNano() <= tmodify.UnixNano() {
		t.Errorf("should set checked when unchanged")
	}
	if status.changed.UnixNano() != tmodify.UnixNano() {
		t.Errorf("should not set changed when unchanged")
	}
}

func TestStatusWeight(t *testing.T) {
	if StatusWeight(StatusOk) > StatusWeight(StatusDegraded) {
		t.Errorf("ok costs less than degraded")
	}

	if StatusWeight(StatusDegraded) > StatusWeight(StatusCritical) {
		t.Errorf("degraded costs less than critical")
	}

	if StatusWeight(StatusCritical) != StatusWeight("Po-taa-tooo") {
		t.Errorf("should default to cost of critical")
	}
}

func TestIsValidStatus(t *testing.T) {
	valids := []string{StatusOk, StatusDegraded, StatusCritical}
	for _, status := range valids {
		if !IsValidStatus(status) {
			t.Errorf("%s is valid status", status)
		}
	}

	if IsValidStatus("Po-taa-tooo") {
		t.Errorf("Po-ta-tooo is not a valid status")
	}
}

func TestParseAndSet(t *testing.T) {
	status := NewServerStatus()

	backend := testutils.NewBackend(0, false)
	defer backend.Shutdown()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", backend.URL()+"/healthz", nil)

	backend.SetStatus(http.StatusOK, "DEGRADED")
	res, _ := client.Do(req)

	status.ParseAndSet(res)
	if status.Current != StatusDegraded {
		t.Errorf("should set status to degraded from header")
	}

	backend.SetStatus(http.StatusInternalServerError, "OK")
	res, _ = client.Do(req)

	status.ParseAndSet(res)
	if status.Current != StatusMaintenance {
		t.Errorf("should set status to maintenance when not ok")
	}
}

func TestCostMasking(t *testing.T) {
	status0 := NewServerStatus()
	status0.Set(StatusDegraded)

	status1 := NewServerStatus()
	status1.Set(StatusCritical)

	if status0.Cost(StatusDegraded) != status1.Cost(StatusCritical) {
		t.Errorf("should mask accepting status")
	}
}

func TestSlowStartFactor(t *testing.T) {
	status := NewServerStatus()

	status.Set(StatusCritical)
	if status.SlowStartFactor() == 0 {
		t.Errorf("should affect critical servers")
	}

	status.Set(StatusDegraded)
	if status.SlowStartFactor() == 0 {
		t.Errorf("should affect degraded servers")
	}

	status.Set(StatusOk)
	fact := status.SlowStartFactor()
	if fact != Kstartup {
		t.Errorf("should be Kstartup")
	}

	time.Sleep(1 * time.Second)
	fact = status.SlowStartFactor()
	if fact >= Kstartup {
		t.Errorf("should decrease monotonically")
	}

	status.changed, _ = time.Parse("1970-Jan-01", "1970-Jan-01")
	fact = status.SlowStartFactor()
	if fact != 0 {
		t.Errorf("should be zero after Kstartup")
	}
}

func TestSlowStartShape(t *testing.T) {
	if !testing.Verbose() {
		t.Skipf("skipping shape test, use verbose to run")
	}

	status := NewServerStatus()
	status.Set(StatusOk)

	t.Logf("Tdelta, SlowStartFactor")
	fact := status.SlowStartFactor()
	for i := 0; i < Tstartup; i++ {
		time.Sleep(1 * time.Second)
		fact_ := status.SlowStartFactor()
		t.Logf("%2d, %d", i, fact_)
		if fact_ < 1 || fact < fact_ {
			t.Errorf("should decrease monotonically to 0")
		}
		fact = fact_
	}

	time.Sleep(100 * time.Millisecond)
	if status.SlowStartFactor() != 0 {
		t.Errorf("should be 0 after Tstartup")
	}
}
