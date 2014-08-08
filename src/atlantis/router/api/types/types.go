
package types 

import (
	cfg "atlantis/router/config
)


type ApiPool struct {
	User	string
	Secret	string
	Pool	cfg.Pool
}

type ApiRule struct {
	User	string
	Secret	string
	Rule	cfg.Rule
}

type ApiTrie struct {
	User	string
	Secret	string
	Trie	cfg.Trie
}

type ApiPort {
	User	string
	Secret	string
	Port	cfg.Port
} 
