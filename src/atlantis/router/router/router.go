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
	poolCb zk.EventCallbacks
	hostCb zk.EventCallbacks
	ruleCb zk.EventCallbacks
	trieCb zk.EventCallbacks
	portCb zk.EventCallbacks

	// configuration
	ZkRoot       string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New(zkServers string) *Router {
	// everything underneath uses atlantis/logger
	logger.InitPkgLogger()

	c := config.NewConfig(routing.DefaultMatcherFactory())
	return &Router{
		zk:     zk.ManagedZkConn(zkServers),
		config: c,
		poolCb: &PoolCallbacks{config: c},
		hostCb: &HostCallbacks{config: c},
		ruleCb: &RuleCallbacks{config: c},
		trieCb: &TrieCallbacks{config: c},
		portCb: nil,

		// configuration
		ZkRoot:       "/atlantis/router",
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
}

func (r *Router) Run() {
	// configuration manager
	go r.reconfigure()

	// launch the status inspector/server
	NewStatusServer(r).Run(8080, 8*time.Second)
}

func (r *Router) reconfigure() {
	zk.SetZkRoot(r.ZkRoot)
	for {
		<-r.zk.ResetCh
		logger.Printf("reloading configuration")
		go r.zk.ManageTree(zk.ZkPaths["pools"], r.poolCb, r.hostCb)
		go r.zk.ManageTree(zk.ZkPaths["rules"], r.ruleCb)
		go r.zk.ManageTree(zk.ZkPaths["tries"], r.trieCb)
		go r.zk.ManageTree(zk.ZkPaths["ports"], r.portCb)
	}
}
