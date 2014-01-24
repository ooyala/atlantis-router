package zk

import (
	"atlantis/router/config"
	"encoding/json"
	"errors"
	"fmt"
	"launchpad.net/gozk"
	"path"
	"strconv"
)

var ZkPaths map[string]string = map[string]string{
	"pools": "/pools",
	"rules": "/rules",
	"tries": "/tries",
	"ports": "/ports",
}

func SetZkRoot(root string) {
	ZkPaths["pools"] = path.Join(root, "pools")
	ZkPaths["rules"] = path.Join(root, "rules")
	ZkPaths["tries"] = path.Join(root, "tries")
	ZkPaths["ports"] = path.Join(root, "ports")
}

func PoolExists(zk *zookeeper.Conn, name string) (bool, error) {
	stat, err := zk.Exists(path.Join(ZkPaths["pools"], name))
	return (stat != nil), err
}

func ListPools(zk *zookeeper.Conn) ([]string, error) {
	pools, _, err := zk.Children(ZkPaths["pools"])
	return pools, err
}

func GetPool(zk *zookeeper.Conn, name string) (config.Pool, error) {
	zkPath := path.Join(ZkPaths["pools"], name)

	jsonBlob, _, err := zk.Get(zkPath)
	if err != nil {
		return config.Pool{}, err
	}

	hosts, err := GetHosts(zk, name)
	if err != nil {
		return config.Pool{}, err
	}

	var zkPool ZkPool

	err = json.Unmarshal([]byte(jsonBlob), &zkPool)
	if err != nil {
		return config.Pool{}, err
	}

	return zkPool.Pool(hosts), nil
}

func SetPool(zk *zookeeper.Conn, pool config.Pool) error {
	zkPath := path.Join(ZkPaths["pools"], pool.Name)

	zkPool, _ := ToZkPool(pool)

	jsonBlob, err := json.Marshal(zkPool)
	if err != nil {
		return err
	}

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}

	if stat == nil {
		_, err = zk.Create(zkPath, string(jsonBlob), 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	} else {
		_, err = zk.Set(zkPath, string(jsonBlob), -1)
	}
	return err
}

func DelPool(zk *zookeeper.Conn, name string) error {
	zkPath := path.Join(ZkPaths["pools"], name)

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}
	if stat == nil {
		return nil
	}

	children, _, err := zk.Children(zkPath)
	if err != nil {
		return err
	}

	for _, child := range children {
		childPath := path.Join(zkPath, child)

		err = zk.Delete(childPath, -1)
		if err != nil {
			return err
		}
	}

	err = zk.Delete(zkPath, -1)
	return err
}

func GetHosts(zk *zookeeper.Conn, pool string) (map[string]config.Host, error) {
	zkPath := path.Join(ZkPaths["pools"], pool)

	children, _, err := zk.Children(zkPath)
	if err != nil {
		return nil, err
	}

	hosts := map[string]config.Host{}

	for _, child := range children {
		childPath := path.Join(zkPath, child)

		jsonBlob, _, err := zk.Get(childPath)
		if err != nil {
			return nil, err
		}

		var host config.Host

		err = json.Unmarshal([]byte(jsonBlob), &host)
		if err != nil {
			return nil, err
		}

		hosts[child] = host
	}

	return hosts, nil
}

func AddHosts(zk *zookeeper.Conn, pool string, hosts map[string]config.Host) error {
	zkPath := path.Join(ZkPaths["pools"], pool)

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}
	if stat == nil {
		return errors.New(zkPath + " does not exist")
	}

	for name, host := range hosts {
		hostPath := path.Join(zkPath, name)
		stat, err := zk.Exists(hostPath)
		if err != nil {
			return err
		}
		if stat != nil {
			continue
		}

		jsonBlob, err := json.Marshal(host)
		if err != nil {
			return err
		}

		_, err = zk.Create(hostPath, string(jsonBlob), 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
		if err != nil {
			return err
		}
	}

	return nil
}

func DelHosts(zk *zookeeper.Conn, pool string, hosts []string) error {
	if len(hosts) == 0 {
		return nil
	}

	zkPath := path.Join(ZkPaths["pools"], pool)

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}
	if stat == nil {
		return errors.New(zkPath + " does not exist")
	}

	zkHosts, err := GetHosts(zk, pool)
	if err != nil {
		return err
	}

	hostsToDel := map[string]bool{}
	for _, host := range hosts {
		hostsToDel[host] = true
	}

	for name, _ := range zkHosts {
		if !hostsToDel[name] {
			// we were given a list of hosts, and
			// this host was not a member of that
			// list...
			continue
		}

		hostPath := path.Join(zkPath, name)

		err := zk.Delete(hostPath, -1)
		if err != nil {
			return err
		}
	}

	return nil
}

func RuleExists(zk *zookeeper.Conn, name string) (bool, error) {
	stat, err := zk.Exists(path.Join(ZkPaths["rules"], name))
	return (stat != nil), err
}

func ListRules(zk *zookeeper.Conn) ([]string, error) {
	rules, _, err := zk.Children(ZkPaths["rules"])
	return rules, err
}

func GetRule(zk *zookeeper.Conn, name string) (rule config.Rule, err error) {
	zkPath := path.Join(ZkPaths["rules"], name)

	jsonBlob, _, err := zk.Get(zkPath)
	if err != nil {
		return config.Rule{}, err
	}

	err = json.Unmarshal([]byte(jsonBlob), &rule)
	return rule, err
}

func SetRule(zk *zookeeper.Conn, rule config.Rule) error {
	zkPath := path.Join(ZkPaths["rules"], rule.Name)

	jsonBlob, err := json.Marshal(rule)
	if err != nil {
		return err
	}

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}

	if stat == nil {
		_, err = zk.Create(zkPath, string(jsonBlob), 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	} else {
		_, err = zk.Set(zkPath, string(jsonBlob), -1)
	}
	return err
}

func DelRule(zk *zookeeper.Conn, name string) error {
	zkPath := path.Join(ZkPaths["rules"], name)

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}
	if stat == nil {
		return nil
	}

	err = zk.Delete(zkPath, -1)
	return err
}

func TrieExists(zk *zookeeper.Conn, name string) (bool, error) {
	stat, err := zk.Exists(path.Join(ZkPaths["tries"], name))
	return (stat != nil), err
}

func ListTries(zk *zookeeper.Conn) ([]string, error) {
	tries, _, err := zk.Children(ZkPaths["tries"])
	return tries, err
}

func GetTrie(zk *zookeeper.Conn, name string) (trie config.Trie, err error) {
	zkPath := path.Join(ZkPaths["tries"], name)

	jsonBlob, _, err := zk.Get(zkPath)
	if err != nil {
		return config.Trie{}, err
	}

	err = json.Unmarshal([]byte(jsonBlob), &trie)
	return trie, err
}

func SetTrie(zk *zookeeper.Conn, trie config.Trie) error {
	zkPath := path.Join(ZkPaths["tries"], trie.Name)

	jsonBlob, err := json.Marshal(trie)
	if err != nil {
		return err
	}

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}

	if stat == nil {
		_, err = zk.Create(zkPath, string(jsonBlob), 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	} else {
		_, err = zk.Set(zkPath, string(jsonBlob), -1)
	}
	return err
}

func DelTrie(zk *zookeeper.Conn, name string) error {
	zkPath := path.Join(ZkPaths["tries"], name)

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}
	if stat == nil {
		return nil
	}

	err = zk.Delete(zkPath, -1)
	return err
}

func PortExists(zk *zookeeper.Conn, port uint16) (bool, error) {
	stat, err := zk.Exists(path.Join(ZkPaths["ports"], fmt.Sprintf("%d", port)))
	return (stat != nil), err
}

func ListPorts(zk *zookeeper.Conn) ([]uint16, error) {
	ports, _, err := zk.Children(ZkPaths["ports"])
	if err != nil {
		return []uint16{}, err
	}
	portUints := make([]uint16, len(ports))
	for i, portStr := range ports {
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return []uint16{}, err
		}
		portUints[i] = uint16(port)
	}
	return portUints, nil
}

func GetPort(zk *zookeeper.Conn, portUint uint16) (port config.Port, err error) {
	zkPath := path.Join(ZkPaths["ports"], fmt.Sprintf("%d", portUint))

	jsonBlob, _, err := zk.Get(zkPath)
	if err != nil {
		return config.Port{}, err
	}

	err = json.Unmarshal([]byte(jsonBlob), &port)
	return port, err
}

func SetPort(zk *zookeeper.Conn, port config.Port) error {
	zkPath := path.Join(ZkPaths["ports"], fmt.Sprintf("%d", port.Port))

	jsonBlob, err := json.Marshal(port)
	if err != nil {
		return err
	}

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}

	if stat == nil {
		_, err = zk.Create(zkPath, string(jsonBlob), 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	} else {
		_, err = zk.Set(zkPath, string(jsonBlob), -1)
	}
	return err
}

func DelPort(zk *zookeeper.Conn, port uint16) error {
	zkPath := path.Join(ZkPaths["ports"], fmt.Sprintf("%d", port))

	stat, err := zk.Exists(zkPath)
	if err != nil {
		return err
	}
	if stat == nil {
		return nil
	}

	err = zk.Delete(zkPath, -1)
	return err
}
