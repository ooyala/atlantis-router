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

package router

import (
	"atlantis/router/config"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Port struct {
	port     uint16
	config   *config.Config
	listener net.Listener
}

func NewPort(p uint16, c *config.Config) (*Port, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", p))
	if err != nil {
		return nil, err
	}
	return &Port{
		port:     p,
		config:   c,
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
