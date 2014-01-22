package router

import (
	"atlantis/router/config"
	"atlantis/router/logger"
	"atlantis/router/routing"
	"atlantis/router/zk"
	"time"
)

type Router struct {
	zk     *zk.ZkConn
	config *config.Config

	// callbacks
	poolCallbacks zk.EventCallbacks
	hostCallbacks zk.EventCallbacks
	ruleCallbacks zk.EventCallbacks
	trieCallbacks zk.EventCallbacks
	portCallbacks zk.EventCallbacks

	// configuration
	ZkRoot       string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New(zkServers string) *Router {
	// all packages use atlantis/logger's global logger
	logger.InitPkgLogger()

	c := config.NewConfig(routing.DefaultMatcherFactory())
	return &Router{
		// zookeeper connection
		ZkRoot: "/atlantis/router",
		zk:     zk.ManagedZkConn(zkServers),

		// configuration management
		config:        c,
		poolCallbacks: &PoolCallbacks{config: c},
		hostCallbacks: &HostCallbacks{config: c},
		ruleCallbacks: &RuleCallbacks{config: c},
		trieCallbacks: &TrieCallbacks{config: c},
		portCallbacks: &PortCallbacks{config: c},

		// global read & write timeouts
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
}

func (r *Router) Run() {
	// configuration manager
	go r.reconfigure()

	// launch the statusz and debug server
	NewStatusServer(r).Run(8080, 8*time.Second)
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
