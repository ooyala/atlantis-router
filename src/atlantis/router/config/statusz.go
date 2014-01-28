package config

import (
	"atlantis/router/logger"
	"encoding/json"
	"fmt"
)

// Serialization expected by the javascript which displays status information, and
// also by services polling /statusz to monitor health of routers and pools.
type StatusZ struct {
	Pool             string `json:"pool"`
	Server           string `json:"server"`
	RequestsInFlight uint32 `json:"requests_in_flight"`
	RequestsServiced uint64 `json:"requests_serviced"`
	Status           string `json:"status"`
	StatusChanged    string `json:"status_changed"`
}

func (c *Config) StatusZJSON() (string, error) {
	var response []StatusZ

	c.RLock()
	for _, pool := range c.Pools {
		for _, server := range pool.Servers {
			s := StatusZ{
				Pool:             pool.Name,
				Server:           server.Address,
				RequestsInFlight: server.Metrics.RequestsInFlight,
				RequestsServiced: server.Metrics.RequestsServiced,
				Status:           server.Status.Current,
				StatusChanged:    fmt.Sprintf("%s", server.Status.Changed),
			}
			response = append(response, s)
		}
	}
	defer c.RUnlock()

	data, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("[statusz json] %s", err)
		return "", err
	}

	return string(data), nil
}
