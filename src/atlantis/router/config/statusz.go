/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

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
