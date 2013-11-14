package lb

import (
	"atlantis/router/config"
	"atlantis/router/routing"
	"atlantis/router/zk"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"
)

type LoadBalancer struct {
	zk     *zk.ZkConn
	config *config.Config

	// callbacks
	poolCb zk.EventCallbacks
	hostCb zk.EventCallbacks
	ruleCb zk.EventCallbacks
	trieCb zk.EventCallbacks

	// configuration
	ZkRoot       string
	Port         uint16
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New(zkServers string) *LoadBalancer {
	c := config.NewConfig(routing.DefaultMatcherFactory())

	return &LoadBalancer{
		zk:     zk.ManagedZkConn(zkServers),
		config: c,
		poolCb: &PoolCallbacks{config: c},
		hostCb: &HostCallbacks{config: c},
		ruleCb: &RuleCallbacks{config: c},
		trieCb: &TrieCallbacks{config: c},

		// configuration
		ZkRoot:       "/atlantis/router",
		Port:         uint16(80),
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
}

func (l *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if pool := l.config.Route(r); pool != nil {
		pool.Handle(w, r)
	} else {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}
}

func (l *LoadBalancer) reconfigure() {
	zk.SetZkRoot(l.ZkRoot)
	for {
		<-l.zk.ResetCh
		log.Println("reloading configuration")
		go l.zk.ManageTree(zk.ZkPaths["pools"], l.poolCb, l.hostCb)
		go l.zk.ManageTree(zk.ZkPaths["rules"], l.ruleCb)
		go l.zk.ManageTree(zk.ZkPaths["tries"], l.trieCb)
	}
}

func (l *LoadBalancer) Run() {
	// configuration manager
	go l.reconfigure()

	server := &http.Server{
		Handler:        l,
		Addr:           fmt.Sprintf(":%d", l.Port),
		ReadTimeout:    l.ReadTimeout,
		WriteTimeout:   l.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("listening on :%d", l.Port)
	panic(server.ListenAndServe())
}

type PoolCallbacks struct {
	config *config.Config
}

func (p *PoolCallbacks) Created(zkPath, jsonBlob string) {
	var zkPool zk.ZkPool
	if err := json.Unmarshal([]byte(jsonBlob), &zkPool); err != nil {
		log.Printf("error unmarshalling pool: %s", err.Error())
		return
	}
	p.config.AddPool(zkPool.Pool(map[string]config.Host{}))
}

func (p *PoolCallbacks) Deleted(zkPath string) {
	p.config.DelPool(path.Base(zkPath))
}

func (p *PoolCallbacks) Changed(path, jsonBlob string) {
	var zkPool zk.ZkPool
	if err := json.Unmarshal([]byte(jsonBlob), &zkPool); err != nil {
		log.Printf("error unmarshalling pool: %s", err.Error())
		return
	}
	p.config.UpdatePool(zkPool.Pool(nil))
}

type HostCallbacks struct {
	config *config.Config
}

func (h *HostCallbacks) splitPath(zkPath string) (string, string) {
	return path.Base(zkPath), path.Base(path.Dir(zkPath))
}

func (h *HostCallbacks) Created(zkPath, jsonBlob string) {
	hostName, poolName := h.splitPath(zkPath)

	var host config.Host
	if err := json.Unmarshal([]byte(jsonBlob), &host); err != nil {
		log.Printf("error unmarshalling host: %s", err.Error())
		return
	}

	if pool := h.config.Pools[poolName]; pool != nil {
		pool.AddServer(hostName, h.config.ConstructServer(host))
	}
}

func (h *HostCallbacks) Deleted(zkPath string) {
	hostName, poolName := h.splitPath(zkPath)
	if pool := h.config.Pools[poolName]; pool != nil {
		pool.DelServer(hostName)
	}
}

func (h *HostCallbacks) Changed(path, jsonBlob string) {
	log.Printf("error: cannot change host %s", path)
}

type RuleCallbacks struct {
	config *config.Config
}

func (p *RuleCallbacks) Created(zkPath, jsonBlob string) {
	var rule config.Rule
	if err := json.Unmarshal([]byte(jsonBlob), &rule); err != nil {
		log.Printf("error unmarshalling rule: %s", err.Error())
		return
	}
	p.config.AddRule(rule)
}

func (p *RuleCallbacks) Deleted(zkPath string) {
	p.config.DelRule(path.Base(zkPath))
}

func (p *RuleCallbacks) Changed(path, jsonBlob string) {
	var rule config.Rule
	if err := json.Unmarshal([]byte(jsonBlob), &rule); err != nil {
		log.Printf("error unmarshalling rule: %s", err.Error())
		return
	}
	p.config.UpdateRule(rule)
}

type TrieCallbacks struct {
	config *config.Config
}

func (p *TrieCallbacks) Created(zkPath, jsonBlob string) {
	var trie config.Trie
	err := json.Unmarshal([]byte(jsonBlob), &trie)
	if err != nil {
		log.Printf("error unmarshalling trie: %s", err.Error())
		return
	}
	p.config.AddTrie(trie)
}

func (p *TrieCallbacks) Deleted(zkPath string) {
	p.config.DelTrie(path.Base(zkPath))
}

func (p *TrieCallbacks) Changed(path, jsonBlob string) {
	var trie config.Trie
	if err := json.Unmarshal([]byte(jsonBlob), &trie); err != nil {
		log.Printf("error unmarshalling trie: %s", err.Error())
		return
	}
	p.config.UpdateTrie(trie)
}
