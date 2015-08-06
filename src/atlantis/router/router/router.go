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
	"atlantis/router/logger"
	"atlantis/router/routing"
	"atlantis/router/zk"
	"time"
)

type Router struct {
	// zookeeper connection
	zk     *zk.ZkConn
	ZkRoot string

	// ports to listen
	ports      map[uint16]*Port
	statusPort uint16

	// configuration management
	config        *config.Config
	poolCallbacks zk.EventCallbacks
	hostCallbacks zk.EventCallbacks
	ruleCallbacks zk.EventCallbacks
	trieCallbacks zk.EventCallbacks
	portCallbacks zk.EventCallbacks

	// configuration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New(zkServers string, statusPort uint16) *Router {
	// all packages use atlantis/logger's global logger
	logger.InitPkgLogger()

	c := config.NewConfig(routing.DefaultMatcherFactory())
	r := &Router{
		ZkRoot: "/atlantis/router",
		zk:     zk.ManagedZkConn(zkServers),

		ports:      map[uint16]*Port{},
		statusPort: statusPort,

		config:        c,
		poolCallbacks: &PoolCallbacks{config: c},
		hostCallbacks: &HostCallbacks{config: c},
		ruleCallbacks: &RuleCallbacks{config: c},
		trieCallbacks: &TrieCallbacks{config: c},

		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	r.portCallbacks = &PortCallbacks{
		router: r,
		config: c,
	}
	return r
}

func (r *Router) Run() {
	// configuration manager
	go r.reconfigure()

	// launch the statusz and debug server
	NewStatusServer(r).Run(r.statusPort, 8*time.Second)
}

func (r *Router) reconfigure() {
	zk.SetZkRoot(r.ZkRoot)
	for {
		<-r.zk.ResetCh
		logger.Printf("reloading configuration")
		go r.zk.ManageTree(zk.ZkPaths["pools"], r.poolCallbacks, r.hostCallbacks)
		go r.zk.ManageTree(zk.ZkPaths["rules"], r.ruleCallbacks)
		go r.zk.ManageTree(zk.ZkPaths["tries"], r.trieCallbacks)
		go r.zk.ManageTree(zk.ZkPaths["ports"], r.portCallbacks)
	}
}

func (r *Router) AddPort(p uint16) {
	port, err := NewPort(p, r.config)
	if err != nil {
		logger.Errorf("%s", err.Error())
		return
	}
	r.ports[p] = port
	go port.Run(r.ReadTimeout, r.WriteTimeout)
}

func (r *Router) DelPort(p uint16) {
	r.ports[p].Shutdown()
}

func (r *Router) IsConnectedToZk() bool {
	return r.zk.IsConnected()
}
