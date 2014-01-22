package frontend

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type Port struct {
	port     uint16
	listener net.Listener
}

func NewPort(p uint16) (*PortListener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%u", p))
	if err != nil {
		return nil, err
	}
	return &PortListener{
		port:     p,
		listener: l,
	}, nil
}

func (p *Port) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("atlantis-arrival-time", fmt.Sprintf("%d", time.Now().UnixNano()))

	if pool := p.config.RoutePort(p.port, r); pool != nil {
		pool.Handle(w, r)
	} else {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}
}

func (p *Port) Run(rout, wout time.Duration) {
	server := http.Server{
		Handler:        p,
		ReadTimeout:    rout,
		WriteTimeout:   wout,
		MaxHeaderBytes: 1 << 20,
	}
	server.Serve(p.listener)
}

func (p *Port) Shutdown() {
	p.listener.Close()
}
