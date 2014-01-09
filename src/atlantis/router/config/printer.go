package config

import (
	"atlantis/router/backend"
	"atlantis/router/logger"
	"atlantis/router/routing"
	"encoding/json"
	"fmt"
	"net/http"
)

// This leaks the abstractions of routing.Trie.Walk() and config.Route()
// and is strictly a debugging aid.
func (c *Config) PrintRouting(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/plain")

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

// Serialization expected by the javascript which displays status information, and
// also by services polling /statusz to monitor health of routers and pools.
type status struct {
	pool               string
	server             string
	requests_in_flight uint32
	requests_serviced  uint64
	status             string
	status_changed     string
}

func (c *Config) PrintStatus(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("accept") == "application/json" {
		var response []status

		c.RLock()
		for _, pool := range c.Pools {
			for _, server := range pool.Servers {
				s := status{
					pool:               pool.Name,
					server:             server.Address,
					requests_in_flight: server.Metrics.RequestsInFlight,
					requests_serviced:  server.Metrics.RequestsServiced,
					status:             server.Status.Current,
					status_changed:     fmt.Sprintf("%s", server.Status.Changed),
				}
				response = append(response, s)
			}
		}
		defer c.RUnlock()

		data, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			logger.Errorf("[config printer %s] %s", r.RemoteAddr, err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(data)
		w.WriteHeader(http.StatusOK)
	} else {
		http.ServeFile(w, r, "statusz.html")
	}
}

func (c *Config) Printer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/statusz" {
		c.PrintStatus(w, r)
	} else {
		c.PrintRouting(w, r)
	}
}
