package config

import (
	"atlantis/router/backend"
	"atlantis/router/routing"
	"fmt"
	"net/http"
)

// This leaks the abstractions of routing.Trie.Walk() and config.Route()
// and is strictly a debugging aid.
func (c *Config) PrintRouting(port uint16, r *http.Request) string {
	c.RLock()
	defer c.RUnlock()

	var next *routing.Trie
	var pool *backend.Pool

	output := fmt.Sprintf("port %d\n", port)
	indent := "  "

	next = c.Ports[port]
	for next != nil {
		if pool != nil {
			output += fmt.Sprintf("%spool %s\n", indent, pool.Name)
			return output
		} else {
			output += fmt.Sprintf("%strie %s\n", indent, next.Name)
		}
		indent += "  "

		for _, rule := range next.List {
			if rule.Dummy {
				output += fmt.Sprintf("%srule %s dummy\n", indent, rule.Name)
			} else if rule.Matcher.Match(r) {
				output += fmt.Sprintf("%srule %s T\n", indent, rule.Name)
				pool, next = rule.PoolPtr, rule.NextPtr
				break
			} else {
				output += fmt.Sprintf("%srule %s F\n", indent, rule.Name)
			}
		}
	}
	return output
}
