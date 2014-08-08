
package zk


import (
	. "atlantis/router/zk"
	cfg "atlantis/router/config"
)

func ListPools() ([]string, error){
	
	return routerzk.ListPools(zkConn)
}

func GetPool(name string) (cfg.Pool, error){
	if name == "" {
		return errors.New("Please specify a name")
	}
	
	return routerzk.GetPool(zkConn, name)
}

func SetPool(pool cfg.Pool) error {

	if pool.Name == "" {
                return errors.New("Please specify a name")
        } else if pool.HealthzEvery == "" {
                return errors.New("Please specify a healthz check frequency")
        } else if pool.HealthzTimeout == "" {
                return errors.New("Please specify a healthz timeout")
        } else if pool.RequestTimeout == "" {
                return errors.New("Please specify a request timeout")
        } // no need to check hosts. an empty pool is still a valid pool


	return routerzk.SetPool(zkConn, pool)
}


func DeletePool(name string) error {
	if name == "" {
		return errors.New("Please specify a name")
	}

	return routerzk.DelPool(zkConn, name)
}
