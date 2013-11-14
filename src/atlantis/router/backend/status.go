package backend

import (
	"net/http"
	"time"
)

const (
	StatusOk          = "OK"
	StatusDegraded    = "DEGRADED"
	StatusCritical    = "CRITICAL"
	StatusMaintenance = "MAINTENANCE"
)

type ServerStatus struct {
	Current string
	checked time.Time
	changed time.Time
}

func NewServerStatus() ServerStatus {
	return ServerStatus{
		Current: StatusMaintenance,
		checked: time.Now(),
		changed: time.Now(),
	}
}

func (s *ServerStatus) Set(status string) {
	s.checked = time.Now()
	if s.Current != status {
		s.Current = status
		s.changed = s.checked
	}

}

func StatusWeight(s string) uint32 {
	switch s {
	case StatusOk:
		return 0x10000000
	case StatusDegraded:
		return 0x30000000
	case StatusCritical:
		return 0x70000000
	default:
		// "CRITICAL".StatusWeight()
		return 0x70000000
	}
}

func IsValidStatus(s string) bool {
	return s == StatusOk || s == StatusDegraded || s == StatusCritical
}

func (s *ServerStatus) ParseAndSet(res *http.Response) {
	if res.StatusCode == http.StatusOK {
		hdr := res.Header.Get("Server-Status")
		if IsValidStatus(hdr) {
			s.Set(hdr)
			return
		}
	}
	s.Set(StatusMaintenance)
}

func (s *ServerStatus) Cost(accept string) uint32 {
	cost := StatusWeight(s.Current) &^ StatusWeight(accept)
	return cost + s.SlowStartFactor()
}

const (
	Tstartup = 60   // Startup time in seconds
	Kstartup = 4096 // Maximum slow start cost
)

func (s *ServerStatus) SlowStartFactor() uint32 {
	if s.Current != StatusOk {
		return 0
	}

	Tdelta := time.Now().Unix() - s.changed.Unix()
	if Tdelta > Tstartup {
		return 0
	} else if Tdelta > 0 {
		k := float64(Kstartup)
		return uint32(k/float64(Tdelta) - k/float64(Tstartup))
	} else {
		// Tdelta is 0
		return Kstartup
	}
}
