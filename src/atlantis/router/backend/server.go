package backend

import (
	"atlantis/router/logger"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Address   string
	Status    ServerStatus
	Metrics   ServerMetrics
	Transport *http.Transport
	copier    *Copier
}

func NewServer(address string) *Server {
	return &Server{
		Address: address,
		Status:  NewServerStatus(),
		Metrics: NewServerMetrics(),
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 32,
		},
		copier: NewCopier(),
	}
}

type ResponseError struct {
	Response *http.Response
	Error    error
}

func (s *Server) RoundTrip(req *http.Request, ch chan ResponseError) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)
			ch <- ResponseError{nil, errors.New(err)}
		}
	}()

	s.Metrics.RequestStart()
	defer s.Metrics.RequestDone()

	req.URL.Scheme = "http"
	req.URL.Host = s.Address

	resp, err := s.Transport.RoundTrip(req)
	if err == nil {
		ch <- ResponseError{resp, nil}
	} else {
		ch <- ResponseError{resp, err}
	}
}

func (s *Server) logPrefix(r *http.Request, tstart time.Time) string {
	now := time.Now()

	var rtt0, rtt1 int64

	// Calculate the total round trip time if front end inserted the atlantis-arrival-time
	// header before routing the request. The header is assumed to be created by calling
	// time.Now().UnixNano() or equivalent.
	arr, err := strconv.ParseInt(r.Header.Get("atlantis-arrival-time"), 10, 64)
	if err != nil {
		rtt0 = now.UnixNano() - arr
	}

	rtt1 = now.UnixNano() - tstart.UnixNano()

	// Log prefix includes server address, request source and URI, and round trip times.
	return fmt.Sprintf("[server %s][request %s|%s][rtt %d|%d]", s.Address, r.RemoteAddr,
		r.URL, rtt0, rtt1)
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request, tout time.Duration) {
	// X-Forwarded-For; we are a proxy.
	ip := strings.Split(r.RemoteAddr, ":")[0]
	r.Header.Add("X-Forwarded-For", ip)

	resErrCh := make(chan ResponseError)
	tstart := time.Now()
	go s.RoundTrip(r, resErrCh)

	select {
	case resErr := <-resErrCh:
		if resErr.Error == nil {
			logger.Printf("%s %d", s.logPrefix(r, tstart), resErr.Response.StatusCode)
			defer resErr.Response.Body.Close()
			for hdr, vals := range resErr.Response.Header {
				for _, val := range vals {
					w.Header().Add(hdr, val)
				}
			}
			w.WriteHeader(resErr.Response.StatusCode)
			_, err := s.copier.Copy(w, resErr.Response.Body)
			if err != nil {
				logger.Errorf("%s %s", s.logPrefix(r, tstart), err)
			}
		} else {
			logger.Errorf("%s %s", s.logPrefix(r, tstart), resErr.Error)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}
	case <-time.After(tout):
		logger.Printf("%s timeout", s.logPrefix(r, tstart))
		http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
	}
}

func (s *Server) CheckStatus(tout time.Duration) {
	r, _ := http.NewRequest("GET", "http://"+s.Address+"/healthz", nil)

	resErrCh := make(chan ResponseError)
	tstart := time.Now()
	go s.RoundTrip(r, resErrCh)

	select {
	case resErr := <-resErrCh:
		if resErr.Error == nil {
			logger.Printf("%s %d", s.logPrefix(r, tstart), resErr.Response.StatusCode)
			defer resErr.Response.Body.Close()
			s.Status.ParseAndSet(resErr.Response)
		} else {
			logger.Errorf("%s %s", s.logPrefix(r, tstart), resErr.Error)
			s.Status.Set(StatusCritical)
		}
	case <-time.After(tout):
		logger.Errorf("%s timeout", s.logPrefix(r, tstart))
		s.Status.Set(StatusCritical)
	}
}

func (s *Server) Cost(accept string) uint32 {
	return s.Status.Cost(accept) + s.Metrics.Cost()
}

func (s *Server) Shutdown() {
	s.copier.Shutdown()
}
