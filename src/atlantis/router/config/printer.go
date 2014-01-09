package config

import(
	"atlantis/router/backend"
	"atlantis/router/routing"
	"fmt"
	"net/http"
)

// This leaks the abstractions of routing.Trie.Walk() and config.Route()
// and is strictly a debugging aid.
func (c *Config) PrintRouting(w http.ResponseWriter, r *http.Request) {
	c.RLock()
	defer c.RUnlock()

	var next *routing.Trie
	var pool *backend.Pool

	var indent string

	next = c.Tries["root"]
	for next != nil {
		if pool != nil {
			fmt.Fprintf(w, "%spool %s", indent, pool.Name)
			return
		} else {
			fmt.Fprintf(w, "%strie %s", indent, next.Name)
		}
		indent += "    "

		for _, rule := range next.List {
			if rule.Dummy {
				fmt.Fprintf(w, "%srule %s dummy", indent, rule.Name)
			}
			if rule.Matcher.Match(r) {
				fmt.Fprintf(w, "%srule %s T", indent, rule.Name)
				pool, next = rule.PoolPtr, rule.NextPtr
				break
			} else {
				fmt.Fprintf(w, "%srule %s F", indent, rule.Name)
			}
		}
	}
}
