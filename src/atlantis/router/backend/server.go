package backend

import (
	"atlantis/router/logger"
	"errors"
	"fmt"
	"net/http"
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

func (s *Server) Handle(w http.ResponseWriter, r *http.Request, tout time.Duration) {
	// X-Forwarded-For; we are a proxy.
	ip := strings.Split(r.RemoteAddr, ":")[0]
	r.Header.Add("X-Forwarded-For", ip)

	resErrCh := make(chan ResponseError)
	go s.RoundTrip(r, resErrCh)

	select {
	case resErr := <-resErrCh:
		if resErr.Error == nil {
			defer resErr.Response.Body.Close()
			for hdr, vals := range resErr.Response.Header {
				for _, val := range vals {
					w.Header().Add(hdr, val)
				}
			}
			w.WriteHeader(resErr.Response.StatusCode)
			_, err := s.copier.Copy(w, resErr.Response.Body)
			if err != nil {
				logger.Errorf("[SERVER %s] %s", s.Address, err)
			}
		} else {
			logger.Errorf("[SERVER %s] %s", s.Address, resErr.Error)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}
	case <-time.After(tout):
		http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
	}
}

func (s *Server) CheckStatus(tout time.Duration) {
	url := fmt.Sprintf("http://%s/healthz", s.Address)
	req, _ := http.NewRequest("GET", url, nil)

	resErrCh := make(chan ResponseError)
	go s.RoundTrip(req, resErrCh)

	select {
	case resErr := <-resErrCh:
		if resErr.Error == nil {
			defer resErr.Response.Body.Close()
			s.Status.ParseAndSet(resErr.Response)
		} else {
			logger.Errorf("[SERVER %s] /healthz %s", resErr.Error)
			s.Status.Set(StatusCritical)
		}
	case <-time.After(tout):
		s.Status.Set(StatusCritical)
	}
}

func (s *Server) Cost(accept string) uint32 {
	return s.Status.Cost(accept) + s.Metrics.Cost()
}

func (s *Server) Shutdown() {
	s.copier.Shutdown()
}
