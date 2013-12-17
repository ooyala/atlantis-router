package lb

import (
	"atlantis/router/config"
	"atlantis/router/logger"
	"atlantis/router/routing"
	"atlantis/router/zk"
	"encoding/json"
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
	ZkRoot              string
	ListenAddr          string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	AtlantisAppSuffixes []string
}

func New(zkServers string) *LoadBalancer {
	c := config.NewConfig(routing.DefaultMatcherFactory())

	logger.InitPkgLogger()

	return &LoadBalancer{
		zk:     zk.ManagedZkConn(zkServers),
		config: c,
		poolCb: &PoolCallbacks{config: c},
		hostCb: &HostCallbacks{config: c},
		ruleCb: &RuleCallbacks{config: c},
		trieCb: &TrieCallbacks{config: c},

		// configuration
		ZkRoot:              "/atlantis/router",
		ListenAddr:          "0.0.0.0:80",
		ReadTimeout:         120 * time.Second,
		WriteTimeout:        120 * time.Second,
		AtlantisAppSuffixes: []string{},
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
	routing.AtlantisAppSuffixes = l.AtlantisAppSuffixes

	// configuration manager
	go l.reconfigure()

	server := &http.Server{
		Handler:        l,
		Addr:           l.ListenAddr,
		ReadTimeout:    l.ReadTimeout,
		WriteTimeout:   l.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Printf("listening on %s", l.ListenAddr)
	panic(server.ListenAndServe())
}

type PoolCallbacks struct {
	config *config.Config
}

func (p *PoolCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("PoolCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	var zkPool zk.ZkPool
	if err := json.Unmarshal([]byte(jsonBlob), &zkPool); err != nil {
		logger.Errorf("%s unmarshalling %s as pool", err.Error(), jsonBlob)
		return
	}
	log.Printf("[config] + pool: %s", zkPool.Name)
	p.config.AddPool(zkPool.Pool(map[string]config.Host{}))
}

func (p *PoolCallbacks) Deleted(zkPath string) {
	logger.Debugf("PoolCallbacks.Deleted(%s)", zkPath)
	p.config.DelPool(path.Base(zkPath))
}

func (p *PoolCallbacks) Changed(path, jsonBlob string) {
	logger.Debugf("PoolCallbacks.Changed(%s, %s)", path, jsonBlob)
	var zkPool zk.ZkPool
	if err := json.Unmarshal([]byte(jsonBlob), &zkPool); err != nil {
		logger.Errorf("%s unmarshalling %s as pool", err.Error(), jsonBlob)
		return
	}
	log.Printf("[config] > pool: %s", zkPool.Name)
	p.config.UpdatePool(zkPool.Pool(nil))
}

type HostCallbacks struct {
	config *config.Config
}

func (h *HostCallbacks) splitPath(zkPath string) (string, string) {
	return path.Base(zkPath), path.Base(path.Dir(zkPath))
}

func (h *HostCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("HostCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	hostName, poolName := h.splitPath(zkPath)

	var host config.Host
	if err := json.Unmarshal([]byte(jsonBlob), &host); err != nil {
		logger.Errorf("%s unmarshalling %s as host", err.Error(), jsonBlob)
		return
	}

	log.Printf("[config] + host: %s %s", poolName, hostName)
	if pool := h.config.Pools[poolName]; pool != nil {
		pool.AddServer(hostName, h.config.ConstructServer(host))
	}
}

func (h *HostCallbacks) Deleted(zkPath string) {
	logger.Debugf("HostCallbacks.Deleted(%s)", zkPath)
	hostName, poolName := h.splitPath(zkPath)
	if pool := h.config.Pools[poolName]; pool != nil {
		pool.DelServer(hostName)
	}
	log.Printf("[config] - host: %s %s", poolName, hostName)
}

func (h *HostCallbacks) Changed(path, jsonBlob string) {
	logger.Errorf("HostCallbacks.Changed(%s, %s)", path, jsonBlob)
}

type RuleCallbacks struct {
	config *config.Config
}

func (p *RuleCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("RuleCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	var rule config.Rule
	if err := json.Unmarshal([]byte(jsonBlob), &rule); err != nil {
		logger.Errorf("%s unmarshalling %s as rule", err.Error(), jsonBlob)
		return
	}
	log.Printf("[config] + rule: %s", rule.Name)
	p.config.AddRule(rule)
}

func (p *RuleCallbacks) Deleted(zkPath string) {
	logger.Debugf("RuleCallbacks.Deleted(%s)", zkPath)
	p.config.DelRule(path.Base(zkPath))
}

func (p *RuleCallbacks) Changed(path, jsonBlob string) {
	logger.Debugf("RuleCallbacks.Changed(%s, %s)", path, jsonBlob)
	var rule config.Rule
	if err := json.Unmarshal([]byte(jsonBlob), &rule); err != nil {
		logger.Errorf("%s unmarshalling %s as rule", err.Error(), jsonBlob)
		return
	}
	log.Printf("[config] > rule: %s", rule.Name)
	p.config.UpdateRule(rule)
}

type TrieCallbacks struct {
	config *config.Config
}

func (p *TrieCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("TrieCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	var trie config.Trie
	err := json.Unmarshal([]byte(jsonBlob), &trie)
	if err != nil {
		logger.Errorf("%s unmarshalling %s as trie", err.Error(), jsonBlob)
		return
	}
	log.Printf("[config] + trie: %s", trie.Name)
	p.config.AddTrie(trie)
}

func (p *TrieCallbacks) Deleted(zkPath string) {
	logger.Debugf("TrieCallbacks.Deleted(%s)", zkPath)
	p.config.DelTrie(path.Base(zkPath))
}

func (p *TrieCallbacks) Changed(path, jsonBlob string) {
	logger.Debugf("TrieCallback.Changed(%s, %s)", path, jsonBlob)
	var trie config.Trie
	if err := json.Unmarshal([]byte(jsonBlob), &trie); err != nil {
		logger.Errorf("%s unmarshalling %s as trie", err.Error(), jsonBlob)
		return
	}
	log.Printf("[config] > trie: %s", trie.Name)
	p.config.UpdateTrie(trie)
}
