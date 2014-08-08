
package zk


import (
	"strconv"
	. "atlantis/router/zk"
	cfg "atlantis/router/config"
)

func GetHosts(poolName string) (map[string]cfg.Host, error) {

	if poolName == "" {
		return nil, errors.New("Please specify a pool name to get the hosts from")
	}

	return routerzk.GetHosts(zkConn, poolName)

}

func AddHosts(poolName string, hosts map[string]cfg.Host) error {

	if poolName == " {
		return errors.New("Please specify a pool name to add the hosts to")
	} else if len(hosts) == 0 {
		return errors.New("Please specify at least one host to add to the pool")
	}

	return routerzk.AddHosts(zkConn, poolName, hosts)
}

func DeleteHosts(poolName string, hosts []string) error {

	if poolName == "" {
		return errors.New("Please specify a pool name to delete hosts from")
	} else if len(hosts) == 0 {
		return errors.New("Please specifiy at least one host to delete from the pool")
	}

	return routerzk.DelHosts(zkConn, poolName, hosts)	
}
