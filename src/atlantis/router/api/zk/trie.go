
package zk


import (
	. "atlantis/router/zk"
	cfg "atlantis/router/config"
)

func ListTries() ([]string, error){
	
	return routerzk.ListTries(zkConn)
}

func GetTrie(name string) (cfg.Trie, error){
	if name == "" {
		return errors.New("Please specify a name")
	}
	
	return routerzk.GetTrie(zkConn, name)
}

func SetTrie(trie cfg.Trie) error {

	if trie.Name == "" {
		return errors.New("Please specify a name")
	} else if len(trie.Rules) <= 0  {
		return errors.New("Please specify a rule")
	}

	return routerzk.SetTrie(zkConn, trie)
}


func DeleteTrie(name string) error {
	if name == "" {
		return errors.New("Please specify a name")
	}

	return routerzk.DelTrie(zkConn, name)
}
