package zk

import (
	"atlantis/router/config"
	"atlantis/router/testutils"
	"encoding/json"
	"launchpad.net/gozk"
	"path"
	"sort"
	"testing"
)

// The sequence of tests in this file assumes that:
//   1. The tests run in the order they are defined.
//   2. Later tests depend on earlier tests passing.

var server *zookeeper.Server
var zkConn *testutils.ZkConn

const zkRoot = "/testing"

// Not a test, but go test won't run it otherwise.
func TestStartServer(t *testing.T) {
	var err error
	server, err = testutils.NewZkServer()
	if err != nil {
		t.Fatalf("cannot start zk server")
	}

	zkConn, err = testutils.NewZkConn(server, true)
	if err != nil {
		t.Fatalf("cannot create connection")
	}

	_, err = zkConn.Conn.Create(zkRoot, "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	if err != nil {
		t.Fatalf("cannot create zk root")
	}

	for _, node := range []string{"pools", "rules", "tries", "ports"} {
		zkPath := path.Join(zkRoot, node)
		_, err = zkConn.Conn.Create(zkPath, "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
		if err != nil {
			t.Fatalf("cannot create node %s", node)
		}
	}

	SetZkRoot(zkRoot)
}

// Not global, so we can read-modify-write in tests.
func testPool() config.Pool {
	return config.Pool{
		Name:     "swimming",
		Internal: true,
		Hosts: map[string]config.Host{
			"test0": config.Host{
				Address: "localhost:8080",
			},
			"test1": config.Host{
				Address: "localhost:8081",
			},
		},
		Config: config.PoolConfig{
			HealthzEvery:   "0s",
			HealthzTimeout: "0s",
			RequestTimeout: "0s",
			Status:         "ITSCOMPLICATED",
		},
	}
}

// Not global, so we can read-modify-write in tests.
func testRule() config.Rule {
	return config.Rule{
		Name:     "ring",
		Type:     "one",
		Value:    "rule",
		Pool:     "all",
		Next:     "doom",
		Internal: false,
	}
}

// Not global, so we can read-modify-write in tests.
func testTrie() config.Trie {
	return config.Trie{
		Name: "colors",
		Rules: []string{
			"violet",
			"indigo",
			"blue",
			"green",
			"yellow",
			"orange",
			"red",
		},
		Internal: false,
	}
}

func testPort() config.Port {
	return config.Port{
		Name: "testport",
		Port: 49152,
		Trie: "testtrie",
	}
}

func TestSetPool(t *testing.T) {
	pool := testPool()

	if err := SetPool(zkConn.Conn, pool); err != nil {
		t.Fatalf("cannot create pool")
	}

	jsonBlob, _, err := zkConn.Conn.Get("/testing/pools/swimming")
	if err != nil {
		t.Fatalf("should create node for pool")
	}

	defer func() {
		if err := zkConn.Conn.Delete("/testing/pools/swimming", -1); err != nil {
			t.Fatalf("coulc not clean up pool node")
		}
	}()

	var node ZkPool

	if err := json.Unmarshal([]byte(jsonBlob), &node); err != nil {
		t.Errorf("should marshal pool correctly")
	} else if node.Name != pool.Name || node.Internal != pool.Internal || node.Config != pool.Config {
		t.Errorf("should create pool accurately")
	}

	hosts, _, err := zkConn.Conn.Children("/testing/pools/swimming")
	if err != nil {
		t.Fatalf("cannot read child nodes")
	} else if len(hosts) != 0 {
		t.Errorf("should not set hosts")
	}

	pool.Config = config.PoolConfig{
		HealthzEvery:   "1s",
		HealthzTimeout: "1s",
		RequestTimeout: "1s",
		Status:         "ITSCOMPLICATED",
	}

	if err := SetPool(zkConn.Conn, pool); err != nil {
		t.Errorf("should modify pool")
		t.Fatalf("skipping ...")
	}

	children, _, err := zkConn.Conn.Children("/testing/pools")
	if err != nil || len(children) != 1 {
		t.Fatalf("should modify node in place")
	}

	jsonBlob, _, err = zkConn.Conn.Get("/testing/pools/swimming")
	if err != nil {
		t.Fatalf("cannot read updated node")
	}

	if err := json.Unmarshal([]byte(jsonBlob), &node); err != nil {
		t.Errorf("should marshal correctly")
	} else if node.Name != pool.Name || node.Internal != pool.Internal || node.Config != pool.Config {
		t.Errorf("should update accurately")
	}
}

func TestAddHosts(t *testing.T) {
	pool := testPool()

	if err := AddHosts(zkConn.Conn, "punting", pool.Hosts); err == nil {
		t.Errorf("should error when pool is not present")
	}

	if err := SetPool(zkConn.Conn, pool); err != nil {
		t.Fatalf("cannot create pool")
	}

	defer func() {
		if err := zkConn.Conn.Delete("/testing/pools/swimming", -1); err != nil {
			t.Fatalf("could not clean up pool node")
		}
	}()

	if err := AddHosts(zkConn.Conn, "swimming", pool.Hosts); err != nil {
		t.Errorf("should add hosts to pool")
	}

	hosts, _, err := zkConn.Conn.Children("/testing/pools/swimming")
	if err != nil {
		t.Fatalf("should create child nodes")
	} else {
		sort.Strings(hosts)
		if len(hosts) != 2 || hosts[0] != "test0" || hosts[1] != "test1" {
			t.Errorf("should set hosts accurately")
		}
	}

	defer func() {
		if err := zkConn.Conn.Delete("/testing/pools/swimming/test0", -1); err != nil {
			t.Fatalf("could not clean up host node")
		}
		if err := zkConn.Conn.Delete("/testing/pools/swimming/test1", -1); err != nil {
			t.Fatalf("could not clean up host node")
		}
	}()

	newHosts := map[string]config.Host{
		"test0": config.Host{
			Address: "localhost:8080",
		},
		"test1": config.Host{
			Address: "localhost:8081",
		},
		"test2": config.Host{
			Address: "localhost:8082",
		},
	}

	if err := AddHosts(zkConn.Conn, "swimming", newHosts); err != nil {
		t.Errorf("should allow adding new hosts")
	}

	defer func() {
		if err := zkConn.Conn.Delete("/testing/pools/swimming/test2", -1); err != nil {
			t.Fatalf("could not clean up host node")
		}
	}()

	hosts, _, err = zkConn.Conn.Children("/testing/pools/swimming")
	if err != nil {
		t.Fatalf("cannot read child nodes")
	} else {
		sort.Strings(hosts)
		if len(hosts) != 3 || hosts[0] != "test0" || hosts[1] != "test1" || hosts[2] != "test2" {
			t.Errorf("should update hosts accurately")
		}
	}

	jsonBlob, _, err := zkConn.Conn.Get("/testing/pools/swimming/test2")
	if err != nil {
		t.Errorf("should create node for added host")
	}

	var node config.Host

	if err := json.Unmarshal([]byte(jsonBlob), &node); err != nil {
		t.Errorf("should marshal correctly")
	} else if node != newHosts["test2"] {
		t.Errorf("should marshal accurately")
	}
}

func TestSetRule(t *testing.T) {
	rule := testRule()

	if err := SetRule(zkConn.Conn, rule); err != nil {
		t.Fatalf("cannot create rule")
	}

	jsonBlob, _, err := zkConn.Conn.Get("/testing/rules/ring")
	if err != nil {
		t.Fatalf("should create rule node")
	}

	defer func() {
		if err := zkConn.Conn.Delete("/testing/rules/ring", -1); err != nil {
			t.Fatalf("could not clean up rule node")
		}
	}()

	var node config.Rule

	err = json.Unmarshal([]byte(jsonBlob), &node)
	if err != nil {
		t.Errorf("should marshal correctly")
	} else if node != rule {
		t.Errorf("should marshal accurately")
	}

	rule.Pool = "orcs"
	if err := SetRule(zkConn.Conn, rule); err != nil {
		t.Fatalf("should modify rule")
	}

	children, _, err := zkConn.Conn.Children("/testing/rules")
	if err != nil || len(children) != 1 {
		t.Errorf("should modify node in place")
	}

	jsonBlob, _, err = zkConn.Conn.Get("/testing/rules/ring")
	if err != nil {
		t.Fatalf("cannot read updated node")
	}

	err = json.Unmarshal([]byte(jsonBlob), &node)
	if err != nil {
		t.Errorf("should marshal correctly")
	} else if node != rule {
		t.Errorf("should update accurately")
	}
}

func TestSetTrie(t *testing.T) {
	trie := testTrie()

	if err := SetTrie(zkConn.Conn, trie); err != nil {
		t.Fatalf("cannot create trie")
	}

	jsonBlob, _, err := zkConn.Conn.Get("/testing/tries/colors")
	if err != nil {
		t.Fatalf("should create trie node")
	}

	defer func() {
		if err := zkConn.Conn.Delete("/testing/tries/colors", -1); err != nil {
			t.Fatalf("could not clean up trie node")
		}
	}()

	var node config.Trie

	err = json.Unmarshal([]byte(jsonBlob), &node)
	if err != nil {
		t.Errorf("should unmarshal correctly")
	} else if node.Name != trie.Name || node.Internal != trie.Internal {
		t.Errorf("should unmarshal accurately")
	}

	trie.Internal = true
	if err := SetTrie(zkConn.Conn, trie); err != nil {
		t.Fatalf("should modify trie")
	}

	children, _, err := zkConn.Conn.Children("/testing/tries")
	if err != nil || len(children) != 1 {
		t.Errorf("should modify node in place")
	}

	jsonBlob, _, err = zkConn.Conn.Get("/testing/tries/colors")
	if err != nil {
		t.Fatalf("cannot read updated node")
	}

	err = json.Unmarshal([]byte(jsonBlob), &node)
	if err != nil {
		t.Errorf("should unmarshal correctly")
	} else if node.Internal != true {
		t.Errorf("should update accurately")
	}
}

func TestSetPort(t *testing.T) {
	port := testPort()

	if err := SetPort(zkConn.Conn, port); err != nil {
		t.Fatalf("cannot create port: %v", err)
	}

	jsonBlob, _, err := zkConn.Conn.Get("/testing/ports/testport")
	if err != nil {
		t.Fatalf("should create port node")
	}

	defer func() {
		if err := zkConn.Conn.Delete("/testing/ports/testport", -1); err != nil {
			t.Fatalf("could not clean up port node")
		}
	}()

	var node config.Port

	err = json.Unmarshal([]byte(jsonBlob), &node)
	if err != nil {
		t.Errorf("should unmarshal correctly")
	} else if node.Name != port.Name || node.Port != port.Port || node.Trie != port.Trie {
		t.Errorf("should unmarshal accurately")
	}

	port.Port = uint16(1337)
	if err := SetPort(zkConn.Conn, port); err != nil {
		t.Fatalf("should modify port")
	}

	children, _, err := zkConn.Conn.Children("/testing/ports")
	if err != nil || len(children) != 1 {
		t.Errorf("should modify node in place")
	}

	jsonBlob, _, err = zkConn.Conn.Get("/testing/ports/testport")
	if err != nil {
		t.Fatalf("cannot read updated node")
	}

	err = json.Unmarshal([]byte(jsonBlob), &node)
	if err != nil {
		t.Errorf("should unmarshal correctly")
	} else if node.Port != uint16(1337) {
		t.Errorf("should update accurately")
	}
}

func TestDelPool(t *testing.T) {
	if err := DelPool(zkConn.Conn, "punting"); err != nil {
		t.Errorf("should silently ignore non existent pool")
	}

	if err := SetPool(zkConn.Conn, testPool()); err != nil {
		t.Fatalf("cannot create pool node")
	}

	if err := AddHosts(zkConn.Conn, "swimming", testPool().Hosts); err != nil {
		t.Fatalf("cannot add hosts")
	}

	if err := DelPool(zkConn.Conn, "swimming"); err != nil {
		t.Fatalf("should delete pool")
	}

	if stat, err := zkConn.Conn.Exists("/testing/pools/swimming"); err != nil {
		t.Fatalf("cannot stat pool")
	} else if stat != nil {
		t.Errorf("should delete pool")
		// Ignore errors, we might have partial removal earlier.
		zkConn.Conn.Delete("/testing/pools/swimming/test0", -1)
		zkConn.Conn.Delete("/testing/pools/swimming/test1", -1)
		if err := zkConn.Conn.Delete("/testing/pools/swimming", -1); err != nil {
			t.Fatalf("could not clean up pool")
		}
	}
}

func TestDelHosts(t *testing.T) {
	if err := DelHosts(zkConn.Conn, "punting", []string{"test0", "test1"}); err == nil {
		t.Errorf("should error on non existent pool")
	}

	if err := SetPool(zkConn.Conn, testPool()); err != nil {
		t.Fatalf("cannot create pool")
	}

	defer func() {
		if err := DelPool(zkConn.Conn, "swimming"); err != nil {
			t.Fatalf("could not clean up pool")
		}
	}()

	if err := AddHosts(zkConn.Conn, "swimming", testPool().Hosts); err != nil {
		t.Fatalf("cannot add hosts")
	}

	if err := DelHosts(zkConn.Conn, "swimming", []string{"test1"}); err != nil {
		t.Errorf("should delete host")
	}

	children, _, err := zkConn.Conn.Children("/testing/pools/swimming")
	if err != nil {
		t.Fatalf("cannot read child nodes")
	} else if len(children) != 1 || children[0] != "test0" {
		t.Errorf("should delete host accurately")
	}

	if err := DelHosts(zkConn.Conn, "swimming", []string{"test0", "test1"}); err != nil {
		t.Fatalf("should ignore non existent hosts")
	}

	children, _, err = zkConn.Conn.Children("/testing/pools/swimming")
	if err != nil {
		t.Fatalf("cannot read child nodes")
	} else if len(children) > 0 {
		t.Errorf("should delete test0")
	}
}

func TestDelRule(t *testing.T) {
	if err := DelRule(zkConn.Conn, "ear-ring"); err != nil {
		t.Errorf("should silently ignore non existent rule")
	}

	if err := SetRule(zkConn.Conn, testRule()); err != nil {
		t.Fatalf("cannot create rule node")
	}

	if err := DelRule(zkConn.Conn, "ring"); err != nil {
		t.Errorf("should delete rule")
	}

	if stat, err := zkConn.Conn.Exists("/testing/rules/ring"); err != nil {
		t.Fatalf("cannot stat rule")
	} else if stat != nil {
		if err := zkConn.Conn.Delete("/testing/rules/ring", -1); err != nil {
			t.Fatalf("could not clean up rule")
		}
	}
}

func TestDelTrie(t *testing.T) {
	if err := DelTrie(zkConn.Conn, "sketches"); err != nil {
		t.Errorf("should silently ignore non existent tries")
	}

	if err := SetTrie(zkConn.Conn, testTrie()); err != nil {
		t.Fatalf("cannot create trie node")
	}

	if err := DelTrie(zkConn.Conn, "colors"); err != nil {
		t.Errorf("should delete trie")
	}

	if stat, err := zkConn.Conn.Exists("/testing/tries/colors"); err != nil {
		t.Fatalf("cannot stat trie")
	} else if stat != nil {
		if err := zkConn.Conn.Delete("/testing/tries/colors", -1); err != nil {
			t.Fatalf("could not clean up trie")
		}
	}
}

func TestDelPort(t *testing.T) {
	if err := DelPort(zkConn.Conn, "sketches"); err != nil {
		t.Errorf("should silently ignore non existent ports")
	}

	if err := SetPort(zkConn.Conn, testPort()); err != nil {
		t.Fatalf("cannot create port node")
	}

	if err := DelPort(zkConn.Conn, "testport"); err != nil {
		t.Errorf("should delete port")
	}

	if stat, err := zkConn.Conn.Exists("/testing/ports/testport"); err != nil {
		t.Fatalf("cannot stat port")
	} else if stat != nil {
		if err := zkConn.Conn.Delete("/testing/ports/testport", -1); err != nil {
			t.Fatalf("could not clean up port")
		}
	}
}

func TestPoolExists(t *testing.T) {
	if exists, err := PoolExists(zkConn.Conn, "swimming"); err != nil {
		t.Fatalf("cannot check if pool exists")
	} else {
		if exists {
			t.Errorf("should be false")
		}
	}

	if err := SetPool(zkConn.Conn, testPool()); err != nil {
		t.Fatalf("cannot create pool")
	}

	defer func() {
		if err := DelPool(zkConn.Conn, "swimming"); err != nil {
			t.Fatalf("could not clean up pool")
		}
	}()

	if exists, err := PoolExists(zkConn.Conn, "swimming"); err != nil {
		t.Fatalf("cannot check if pool exists")
	} else {
		if !exists {
			t.Errorf("should be true")
		}
	}
}

func TestRuleExists(t *testing.T) {
	if exists, err := RuleExists(zkConn.Conn, "ring"); err != nil {
		t.Fatalf("cannot check if rule exists")
	} else {
		if exists {
			t.Errorf("should be false")
		}
	}

	if err := SetRule(zkConn.Conn, testRule()); err != nil {
		t.Fatalf("cannot create rule")
	}

	defer func() {
		if err := DelRule(zkConn.Conn, "ring"); err != nil {
			t.Fatalf("could not clean up rule")
		}
	}()

	if exists, err := RuleExists(zkConn.Conn, "ring"); err != nil {
		t.Fatalf("cannot check if rule exists")
	} else {
		if !exists {
			t.Errorf("should be true")
		}
	}
}

func TestTrieExists(t *testing.T) {
	if exists, err := TrieExists(zkConn.Conn, "colors"); err != nil {
		t.Fatalf("cannot check if trie exists")
	} else {
		if exists {
			t.Errorf("should be false")
		}
	}

	if err := SetTrie(zkConn.Conn, testTrie()); err != nil {
		t.Fatalf("cannot create trie")
	}

	defer func() {
		if err := DelTrie(zkConn.Conn, "colors"); err != nil {
			t.Fatalf("could not clean up trie")
		}
	}()

	if exists, err := TrieExists(zkConn.Conn, "colors"); err != nil {
		t.Fatalf("cannot check if trie exists")
	} else {
		if !exists {
			t.Errorf("should be true")
		}
	}
}

func TestPortExists(t *testing.T) {
	if exists, err := PortExists(zkConn.Conn, "testport"); err != nil {
		t.Fatalf("cannot check if port exists")
	} else {
		if exists {
			t.Errorf("should be false")
		}
	}

	if err := SetPort(zkConn.Conn, testPort()); err != nil {
		t.Fatalf("cannot create port")
	}

	defer func() {
		if err := DelPort(zkConn.Conn, "testport"); err != nil {
			t.Fatalf("could not clean up port")
		}
	}()

	if exists, err := PortExists(zkConn.Conn, "testport"); err != nil {
		t.Fatalf("cannot check if port exists")
	} else {
		if !exists {
			t.Errorf("should be true")
		}
	}
}

func TestListPools(t *testing.T) {
	list, err := ListPools(zkConn.Conn)
	if err != nil || len(list) > 0 {
		t.Errorf("should be empty")
	}

	pool := testPool()

	if err := SetPool(zkConn.Conn, pool); err != nil {
		t.Fatalf("cannot set pool")
	}

	defer func() {
		if err := DelPool(zkConn.Conn, "swimming"); err != nil {
			t.Fatalf("could not clean up pool")
		}
	}()

	list, err = ListPools(zkConn.Conn)
	if err != nil {
		t.Errorf("should list pools")
	} else if len(list) != 1 || list[0] != "swimming" {
		t.Errorf("should list pools accurately")
	}

	pool.Name = "punting"

	if err := SetPool(zkConn.Conn, pool); err != nil {
		t.Fatalf("cannot set pool")
	}

	defer func() {
		if err := DelPool(zkConn.Conn, "punting"); err != nil {
			t.Fatalf("could not clean up pool")
		}
	}()

	list, err = ListPools(zkConn.Conn)
	if err != nil {
		t.Errorf("should list pools")
	} else {
		sort.Strings(list)
		if len(list) != 2 || list[0] != "punting" || list[1] != "swimming" {
			t.Errorf("should list pools accurately")
		}
	}
}

func TestListRules(t *testing.T) {
	list, err := ListRules(zkConn.Conn)
	if err != nil || len(list) > 0 {
		t.Errorf("should be empty")
	}

	rule := testRule()

	if err := SetRule(zkConn.Conn, rule); err != nil {
		t.Fatalf("cannot set rule")
	}

	defer func() {
		if err := DelRule(zkConn.Conn, "ring"); err != nil {
			t.Fatalf("could not clean up rule")
		}
	}()

	list, err = ListRules(zkConn.Conn)
	if err != nil {
		t.Errorf("should list rules")
	} else if len(list) != 1 || list[0] != "ring" {
		t.Errorf("should list rules accurately")
	}

	rule.Name = "ear-ring"

	if err := SetRule(zkConn.Conn, rule); err != nil {
		t.Fatalf("cannot set rule")
	}

	defer func() {
		if err := DelRule(zkConn.Conn, "ear-ring"); err != nil {
			t.Fatalf("could not clean up rule")
		}
	}()

	list, err = ListRules(zkConn.Conn)
	if err != nil {
		t.Errorf("should list rules")
	} else {
		sort.Strings(list)
		if len(list) != 2 || list[0] != "ear-ring" || list[1] != "ring" {
			t.Errorf("should list rules accurately")
		}
	}
}

func TestListTries(t *testing.T) {
	list, err := ListTries(zkConn.Conn)
	if err != nil || len(list) > 0 {
		t.Errorf("should be empty")
	}

	trie := testTrie()

	if err := SetTrie(zkConn.Conn, trie); err != nil {
		t.Fatalf("cannot set trie")
	}

	defer func() {
		if err := DelTrie(zkConn.Conn, "colors"); err != nil {
			t.Fatalf("could not clean up trie")
		}
	}()

	list, err = ListTries(zkConn.Conn)
	if err != nil {
		t.Errorf("should list tries")
	} else if len(list) != 1 || list[0] != "colors" {
		t.Errorf("should list tries accurately")
	}

	trie.Name = "sketches"

	if err := SetTrie(zkConn.Conn, trie); err != nil {
		t.Fatalf("could not set trie")
	}

	defer func() {
		if err := DelTrie(zkConn.Conn, "sketches"); err != nil {
			t.Fatalf("could not clean up trie")
		}
	}()

	list, err = ListTries(zkConn.Conn)
	if err != nil {
		t.Errorf("should list tries")
	} else {
		sort.Strings(list)
		if len(list) != 2 || list[0] != "colors" || list[1] != "sketches" {
			t.Errorf("should list tries accurately")
		}
	}
}

func TestListPorts(t *testing.T) {
	list, err := ListPorts(zkConn.Conn)
	if err != nil || len(list) > 0 {
		t.Errorf("should be empty")
	}

	port := testPort()

	if err := SetPort(zkConn.Conn, port); err != nil {
		t.Fatalf("cannot set port")
	}

	defer func() {
		if err := DelPort(zkConn.Conn, "testport"); err != nil {
			t.Fatalf("could not clean up port")
		}
	}()

	list, err = ListPorts(zkConn.Conn)
	if err != nil {
		t.Errorf("should list ports")
	} else if len(list) != 1 || list[0] != "testport" {
		t.Errorf("should list ports accurately")
	}

	port.Name = "testport2"

	if err := SetPort(zkConn.Conn, port); err != nil {
		t.Fatalf("could not set port")
	}

	defer func() {
		if err := DelPort(zkConn.Conn, "testport2"); err != nil {
			t.Fatalf("could not clean up port")
		}
	}()

	list, err = ListPorts(zkConn.Conn)
	if err != nil {
		t.Errorf("should list ports")
	} else {
		sort.Strings(list)
		if len(list) != 2 || list[0] != "testport" || list[1] != "testport2" {
			t.Errorf("should list ports accurately")
		}
	}
}

func TestGetPool(t *testing.T) {
	if _, err := GetPool(zkConn.Conn, "punting"); err == nil {
		t.Errorf("should error on non existent pool")
	}

	pool := testPool()

	if err := SetPool(zkConn.Conn, pool); err != nil {
		t.Fatalf("could not set pool")
	}

	defer func() {
		if err := DelPool(zkConn.Conn, "swimming"); err != nil {
			t.Fatalf("could not clean up pool")
		}
	}()

	node, err := GetPool(zkConn.Conn, "swimming")
	if err != nil {
		t.Errorf("should get pool")
	} else if len(node.Hosts) > 0 {
		t.Errorf("should be empty")
	}

	if err := AddHosts(zkConn.Conn, "swimming", pool.Hosts); err != nil {
		t.Fatalf("could not add hosts to pool")
	}

	node, err = GetPool(zkConn.Conn, "swimming")
	if err != nil {
		t.Errorf("should get pool")
	} else {
		if node.Name != pool.Name || node.Internal != pool.Internal ||
			node.Config != pool.Config || node.Hosts["test0"] != pool.Hosts["test0"] ||
			node.Hosts["test1"] != pool.Hosts["test1"] {
			t.Errorf("should get pool accurately")
		}
	}
}

func TestGetRule(t *testing.T) {
	if _, err := GetRule(zkConn.Conn, "ear-ring"); err == nil {
		t.Errorf("should error on non existent rule")
	}

	rule := testRule()

	if err := SetRule(zkConn.Conn, rule); err != nil {
		t.Fatalf("could not set rule")
	}

	defer func() {
		if err := DelRule(zkConn.Conn, "ring"); err != nil {
			t.Fatalf("could not clean up rule")
		}
	}()

	node, err := GetRule(zkConn.Conn, "ring")
	if err != nil {
		t.Errorf("should get rule")
	} else if node != rule {
		t.Errorf("should get rule accurately")
	}
}

func TestGetTrie(t *testing.T) {
	if _, err := GetTrie(zkConn.Conn, "sketches"); err == nil {
		t.Errorf("should error on non existent trie")
	}

	trie := testTrie()

	if err := SetTrie(zkConn.Conn, trie); err != nil {
		t.Fatalf("could not set trie")
	}

	defer func() {
		if err := DelTrie(zkConn.Conn, "colors"); err != nil {
			t.Fatalf("could not clean up trie")
		}
	}()

	node, err := GetTrie(zkConn.Conn, "colors")
	if err != nil {
		t.Errorf("should get trie")
	} else {
		if node.Name != trie.Name || node.Internal != trie.Internal {
			t.Errorf("should get trie accurately")
		}
		sort.Strings(node.Rules)
		sort.Strings(trie.Rules)
		for i, _ := range trie.Rules {
			if node.Rules[i] != trie.Rules[i] {
				t.Errorf("should get trie accurately")
			}
		}
	}
}

func TestGetPort(t *testing.T) {
	if _, err := GetPort(zkConn.Conn, "sketches"); err == nil {
		t.Errorf("should error on non existent port")
	}

	port := testPort()

	if err := SetPort(zkConn.Conn, port); err != nil {
		t.Fatalf("could not set port")
	}

	defer func() {
		if err := DelPort(zkConn.Conn, "testport"); err != nil {
			t.Fatalf("could not clean up port")
		}
	}()

	node, err := GetPort(zkConn.Conn, "testport")
	if err != nil {
		t.Errorf("should get port")
	} else {
		if node.Name != port.Name || node.Port != port.Port || node.Trie != port.Trie {
			t.Errorf("should get port accurately")
		}
	}
}

// Not a test, but go test won't run it otherwise.
func TestStopServer(t *testing.T) {
	zkConn.Conn.Close()
	server.Destroy()
}
