/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package backend

import (
	"atlantis/router/logger"
	"net/http"
	"strings"
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
	Headers map[string]string //response headers set in case apps are unreachable
	Servers map[string]*Server
	Config  PoolConfig
	killCh  chan bool
	Metrics ConnectionMetrics
}

func DummyPool(name string) *Pool {
	return &Pool{
		Name:  name,
		Dummy: true,
	}
}

func NewPool(name string, config PoolConfig, headers map[string]string) *Pool {
	pool := &Pool{
		Name:    name,
		Dummy:   false,
		Headers: headers,
		Servers: map[string]*Server{},
		Config:  config,
		killCh:  make(chan bool),
		Metrics: NewConnectionMetrics(),
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
		logger.Errorf("[pool %s] server %s exists", p.Name, name)
		return
	}
	p.Servers[name] = server
}

func (p *Pool) DelServer(name string) {
	if _, ok := p.Servers[name]; !ok {
		logger.Errorf("[pool %s] server %s absent", p.Name, name)
		return
	}

	delete(p.Servers, name)
}

func (p *Pool) Reconfigure(config PoolConfig, headers map[string]string) {
	p.Config = config
	p.Headers = headers
}

func (p *Pool) RunChecks() {
	for {
		select {
		case <-time.After(p.Config.HealthzEvery):
			for _, server := range p.Servers {
				go server.CheckStatus(p.Config.HealthzTimeout)
			}
		case <-p.killCh:
			logger.Debugf("[pool %s] stopping checks", p.Name)
			return
		}
	}
}

func (p Pool) Next() *Server {
	var next *Server
	var cost uint32 = 0xffffffff

	for _, server := range p.Servers {
		// Never send traffic to servers under maintenance or unknown.
		if strings.EqualFold(server.Status.Current, StatusMaintenance) || strings.EqualFold(server.Status.Current, StatusUnknown) {
			continue
		}

		newCost := server.Cost(p.Config.Status)
		if newCost < cost {
			next, cost = server, newCost
		}
	}

	return next
}
func (p *Pool) Handle(logRecord *logger.HAProxyLogRecord) {
	pTime := time.Now()
	if p.Dummy {
		logger.Printf("[pool %s] Dummy", p.Name)
		logRecord.Error(logger.BadGatewayMsg, http.StatusBadGateway)
		logRecord.Terminate("Pool: " + logger.BadGatewayMsg)
		return
	}
	p.Metrics.ConnectionStart()
	defer p.Metrics.ConnectionDone()

	server := p.Next()
	if server == nil {
		for k, v := range p.Headers {
			logRecord.AddResponseHeader(k, v)
		}
		// reachable when all servers in pool report StatusMaintenance
		logger.Printf("[pool %s] no server", p.Name)
		logRecord.Error(logger.ServiceUnavailableMsg, http.StatusServiceUnavailable)
		logRecord.Terminate("Pool: " + logger.ServiceUnavailableMsg)
		return
	}
	logRecord.PoolUpdateRecord(p.Name, p.Metrics.GetActiveConnections(), p.Metrics.GetTotalConnections(), pTime)
	server.Handle(logRecord, p.Config.RequestTimeout, &p.Headers)
}
