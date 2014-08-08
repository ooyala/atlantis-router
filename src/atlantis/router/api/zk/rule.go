
package zk


import (
	. "atlantis/router/zk"
	cfg "atlantis/router/config"
)

func ListRules() ([]string, error){
	
	return routerzk.ListRules(zkConn)
}

func GetRule(name string) (cfg.Rule, error){
	if name == "" {
		return errors.New("Please specify a name")
	}
	
	return routerzk.GetRule(zkConn, name)
}

func SetRule(rule cfg.Rule) error {

	if rule.Name == "" {
		return errors.New("Please specify a name")
	} else if rule.Type == "" {
		return errors.New("Please specify a type")
	} else if rule.Value == "" {
		return erros.New("Please specify a value")
	} else if rule.Next == "" {
		return errors.New("Please specify a next value")
	} else if rule.Pool == "" {
		return errors.New("Please specify a pool")
	}

	return routerzk.SetRule(zkConn, rule)
}


func DeleteRule(name string) error {
	if name == "" {
		return errors.New("Please specify a name")
	}

	return routerzk.DelRule(zkConn, name)
}
