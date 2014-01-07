package backend

import (
	"atlantis/router/logger"
	"net/http"
	"time"
)

type PoolConfig struct {
	HealthzEvery   time.Duration
	HealthzTimeout time.Duration
	RequestTimeout time.Duration
	Status         string
}

type Pool struct {
	Name    string
	Dummy   bool
	Servers map[string]*Server
	Config  PoolConfig
	killCh  chan bool
}

func DummyPool(name string) *Pool {
	return &Pool{
		Name:  name,
		Dummy: true,
	}
}

func NewPool(name string, config PoolConfig) *Pool {
	pool := &Pool{
		Name:    name,
		Dummy:   false,
		Servers: map[string]*Server{},
		Config:  config,
		killCh:  make(chan bool),
	}

	go pool.RunChecks()

	return pool
}

func (p *Pool) Shutdown() {
	if !p.Dummy {
		p.killCh <- true
	}
}

func (p *Pool) AddServer(name string, server *Server) {
	if _, ok := p.Servers[name]; ok {
		logger.Errorf("[POOL %s] server %s exists", p.Name, name)
		return
	}
	p.Servers[name] = server
}

func (p *Pool) DelServer(name string) {
	if _, ok := p.Servers[name]; !ok {
		logger.Errorf("[POOL %s] server %s absent", p.Name, name)
		return
	}
	p.Servers[name].Shutdown()
	delete(p.Servers, name)
}

func (p *Pool) Reconfigure(config PoolConfig) {
	p.Config = config
}

func (p *Pool) RunChecks() {
	for {
		select {
		case <-time.After(p.Config.HealthzEvery):
			for _, server := range p.Servers {
				go server.CheckStatus(p.Config.HealthzTimeout)
			}
		case <-p.killCh:
			return
		}
	}
}

func (p Pool) Next() *Server {
	var next *Server
	var cost uint32 = 0xffffffff

	for _, server := range p.Servers {
		// Never send traffic to servers under maintenance.
		if server.Status.Current == StatusMaintenance {
			continue
		}

		newCost := server.Cost(p.Config.Status)
		if newCost < cost {
			next, cost = server, newCost
		}
	}

	return next
}

func (p *Pool) Handle(w http.ResponseWriter, r *http.Request) {
	if p.Dummy {
		logger.Printf("[POOL %s] Dummy", p.Name)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	server := p.Next()
	if server == nil {
		// reachable when all servers in pool report StatusMaintenance
		logger.Errorf("[POOL %s] no server")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	server.Handle(w, r, p.Config.RequestTimeout)
}
