
package zk


import (
	"strconv"
	. "atlantis/router/zk"
	cfg "atlantis/router/config"
)

func ListPorts() ([]uint16, error){
	
	return routerzk.ListPorts(zkConn)
}

func GetPort(name string) (cfg.Port, error){
	if name == "" {
		return errors.New("Please specify a port")
	}

	pUint, err := strconv.ParseUint(name, 10, 16)

	if err != nil {
		return err
	}	
	
	return routerzk.GetPort(zkConn, uint16(pUint))
}

func SetPort(port cfg.Port) error {

	if port.Port == 0 {
		return errors.New("Please specify a port")
	} else if port.Trie == ""  {
		return errors.New("Please specify a trie")
	}

	return routerzk.SetPort(zkConn, port)
}


func DeletePort(name string) error {
	if name == "" {
		return errors.New("Please specify a port")
	}

	pUint, err := strconv.ParseUint(name, 10, 16)
	
	if err != nil {
		return err
	}

	return routerzk.DelPort(zkConn, uint16(pUint))
}
