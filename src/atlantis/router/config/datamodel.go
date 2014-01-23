package config

import (
	"fmt"
)

type PoolConfig struct {
	HealthzEvery   string
	HealthzTimeout string
	RequestTimeout string
	Status         string
}

func (p PoolConfig) Equals(o PoolConfig) bool {
	return p.HealthzEvery == o.HealthzEvery && p.HealthzTimeout == o.HealthzTimeout &&
		p.RequestTimeout == o.RequestTimeout && p.Status == o.Status
}

func (p PoolConfig) StringIndent(i string) (str string) {
	str += fmt.Sprintf("%s--Pool Configuration\n", i)
	str += fmt.Sprintf("%s  Healthz Every   : %s\n", i, p.HealthzEvery)
	str += fmt.Sprintf("%s  Healthz Timeout : %s\n", i, p.HealthzTimeout)
	str += fmt.Sprintf("%s  Request Timeout : %s\n", i, p.RequestTimeout)
	str += fmt.Sprintf("%s  Status          : %s\n", i, p.Status)
	return
}

func (p PoolConfig) String() string {
	return p.StringIndent("")
}

type Host struct {
	Address string
}

func (h Host) Equals(o Host) bool {
	return h.Address == o.Address
}

func (h Host) StringIndent(i string) (str string) {
	str += fmt.Sprintf("%s--Host\n", i)
	str += fmt.Sprintf("%s  Address : %s\n", i, h.Address)
	return
}

func (h Host) String() string {
	return h.StringIndent("")
}

type Pool struct {
	Name     string
	Internal bool
	Hosts    map[string]Host
	Config   PoolConfig
}

func (p Pool) Equals(o Pool) bool {
	return p.Name == o.Name
}

func (p Pool) StringIndent(i string) (str string) {
	str += fmt.Sprintf("%s--Pool\n", i)
	str += fmt.Sprintf("%s  Name     : %s\n", i, p.Name)
	str += fmt.Sprintf("%s  Internal : %t\n", i, p.Internal)
	str += fmt.Sprintf("%s  --Hosts\n", i)
	for name, host := range p.Hosts {
		str += fmt.Sprintf("%s    %s : %s\n", i, name, host.Address)
	}
	str += p.Config.StringIndent(i + "  ")
	return
}

func (p Pool) String() string {
	return p.StringIndent("")
}

type Rule struct {
	Name     string
	Type     string
	Value    string
	Next     string
	Pool     string
	Internal bool
}

func (r Rule) Equals(o Rule) bool {
	return r.Name == o.Name
}

func (r Rule) StringIndent(i string) (str string) {
	str += fmt.Sprintf("%s--Rule\n", i)
	str += fmt.Sprintf("%s  Name     : %s\n", i, r.Name)
	str += fmt.Sprintf("%s  Internal : %t\n", i, r.Internal)
	str += fmt.Sprintf("%s  Type     : %s\n", i, r.Type)
	str += fmt.Sprintf("%s  Value    : %s\n", i, r.Value)
	str += fmt.Sprintf("%s  Next     : %s\n", i, r.Next)
	str += fmt.Sprintf("%s  Pool     : %s\n", i, r.Pool)
	return
}

func (r Rule) String() string {
	return r.StringIndent("")
}

type Trie struct {
	Name     string
	Rules    []string
	Internal bool
}

func (t Trie) Equals(o Trie) bool {
	return t.Name == o.Name
}

func (t Trie) StringIndent(i string) (str string) {
	str += fmt.Sprintf("%s--Trie\n", i)
	str += fmt.Sprintf("%s  Name     : %s\n", i, t.Name)
	str += fmt.Sprintf("%s  Internal : %t\n", i, t.Internal)
	str += fmt.Sprintf("%s  --Rules\n", i)
	for _, rule := range t.Rules {
		str += fmt.Sprintf("%s    Rule : %s\n", i, rule)
	}
	return
}

func (t *Trie) String() string {
	return t.StringIndent("")
}

type Port struct {
	Port     uint16
	Trie     string
	Internal bool
}

func (p Port) Equals(o Port) bool {
	return p.Port == o.Port
}

func (p Port) StringIndent(i string) (str string) {
	str += fmt.Sprintf("%s--Port\n", i)
	str += fmt.Sprintf("%s  Internal : %t\n", i, p.Internal)
	str += fmt.Sprintf("%s  Port     : %d\n", i, p.Port)
	str += fmt.Sprintf("%s  Trie     : %s\n", i, p.Trie)
	return
}

func (p *Port) String() string {
	return p.StringIndent("")
}
