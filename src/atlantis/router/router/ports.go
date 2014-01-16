package router

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type PortListener struct {
	router   *Router
	trie     string
	listener net.Listener
}

func NewPortListener(p uint16, r *Router, t string) (*PortListener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%u", p))
	if err != nil {
		return nil, err
	}
	return &PortListener{
		router:   r,
		listener: l,
	}, nil
}

func (p *PortListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("atlantis-arrival-time", fmt.Sprintf("%d", time.Now().UnixNano()))

	if pool := p.router.config.RouteFrom(p.trie, r); pool != nil {
		pool.Handle(w, r)
	} else {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}
}

func (p *PortListener) Run(r_tout, w_tout time.Duration) {
	server := http.Server{
		Handler:        p,
		ReadTimeout:    r_tout,
		WriteTimeout:   w_tout,
		MaxHeaderBytes: 1 << 20,
	}
	server.Serve(p.listener)
}

func (p *PortListener) Shutdown() {
	p.listener.Close()
}
